package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/xanzy/go-gitlab"
)

type APIToken string
type Records [][]string

type BaseService struct {
	exporter        *Exporter
	filename        string
	attachmentField *string
}

type Exporter struct {
	Client *gitlab.Client `json:"client"`

	State *State `json:"state"`
	//	State map[string]interface{}

	CurrentUser *User `json:"user"`

	Exports Records  `json:"exports"`
	BaseAPI *url.URL `json:"base_api"`
	Token   APIToken `json:"token"`

	BaseDir          string `json:"base_dir"`
	TmpAttachmentDir string `json:"tmp_attachment_dir"`
	TmpRepositoryDir string `json:"tmp_repository_dir"`
	StaticDir        string `json:"static_dir"`
	TmpDir           string `json:"tmp_dir"`

	AttachmentRegex *regexp.Regexp `json:"attachment_regex"`

	CurrentProject *gitlab.Project `json:"current_project"`

	Attachments    *AttachmentService    `json:"attachments"`
	Branches       *BranchService        `json:"branches"`
	CommitComments *CommitCommentService `json:"commit_comments"`
	Groups         *GroupService         `json:"groups"`
	Issues         *IssueService         `json:"issues"`
	Labels         *LabelService         `json:"labels"`
	MergeRequests  *MergeRequestService  `json:"merge_requests"`
	Milestones     *MilestoneService     `json:"milestones"`
	Projects       *ProjectService       `json:"projects"`
	Repositories   *RepositoryService    `json:"repository"`
	Tags           *TagService           `json:"tags"`
	Users          *UserService          `json:"users"`
	Webhooks       *WebhookService       `json:"webhooks"`
}

type State struct {
	Attachments    []Attachment    `json:"attachments"`
	Branches       []Branch        `json:"branches"`
	Collaborators  []User          `json:"collaborators"`
	CommitComments []CommitComment `json:"commit_comments"`
	Issues         []Issue         `json:"issues"`
	Groups         []Group         `json:"groups"`
	Labels         []Label         `json:"labels"`
	MergeRequests  []MergeRequest  `json:"merge_requests"`
	Milestones     []Milestone     `json:"milestones"`
	Projects       []Project       `json:"projects"`
	Repositories   []Repository    `json:"repositories"`
	Tags           []Tag           `json:"tags"`
	Users          []User          `json:"users"`
	Webhooks       []Webhook       `json:"webhooks"`
}

func NewExporter(exports Records, baseAPI string) *Exporter {
	uri, err := url.Parse(baseAPI)
	if err != nil {
		log.Fatalf("Failed to parse the given REST API URL: %v", err)
	}

	apiToken, isSet := os.LookupEnv("GITLAB_API_PRIVATE_TOKEN")
	if apiToken == "" || !isSet {
		panic("[ERROR] Must set $GITLAB_API_PRIVATE_TOKEN!")
	}

	gitlabClient, err := gitlab.NewClient(apiToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	//	type State StateMap

	//	state := map[string]interface{}{
	//		"Attachments":    []Attachment{},
	//		"CommitComments": []CommitComment{},
	//		"Issues":         []Issue{},
	//		"Labels":         []Label{},
	//		"MergeRequests":  []MergeRequest{},
	//		"Milestones":     []Milestone{},
	//		"Projects":       []Project{},
	//		"Repositories":   []Repository{},
	//		"Tags":           []Tag{},
	//		"Users":          []User{},
	//		"Webhooks":       []Webhook{},
	//	}

	e := &Exporter{
		Client: gitlabClient,
		State: &State{ // Set defalt values so "null" is never a value in the json.
			Collaborators: make([]User, 0),
			Labels:        make([]Label, 0),
			Webhooks:      make([]Webhook, 0),
		},
		Exports:         exports,
		BaseAPI:         uri,
		BaseDir:         "migration",
		StaticDir:       "static",
		Token:           APIToken(apiToken),
		AttachmentRegex: regexp.MustCompile(`^\[(.*)\]\((.*)\)$`),
	}

	e.CreateTmpDirs()

	err = e.WriteStaticFiles()
	if err != nil {
		log.Fatalf("Failed to read entries in directory `%s`: %v", e.StaticDir, err)
	}

	e.Attachments = NewAttachmentService(e)
	e.Branches = NewBranchService(e)
	e.CommitComments = NewCommitCommentService(e)
	e.Groups = NewGroupService(e)
	e.Issues = NewIssueService(e)
	e.Labels = NewLabelService(e)
//	e.MergeRequests = NewMergeRequestService(e)
	e.Milestones = NewMilestoneService(e)
	e.Projects = NewProjectService(e)
	e.Repositories = NewRepositoryService(e)
	e.Tags = NewTagService(e)
	e.Users = NewUserService(e)
	e.Webhooks = NewWebhookService(e)

	user, err := e.Users.GetUser()
	if err != nil {
		log.Fatalf("Could not get details for current user: %v", err)
	}
	e.CurrentUser = user

	return e
}

// Taken from `https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726`.
func (e *Exporter) Compress() error {
	src := e.TmpDir

	filename := "migration.tar.gz"
	buf, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not get create filename `%s`: %v", filename, err)
	}

	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// is file a folder?
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	mode := fi.Mode()
	if mode.IsRegular() {
		// get header
		header, err := tar.FileInfoHeader(fi, src)
		if err != nil {
			return err
		}
		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// get content
		data, err := os.Open(src)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}
	} else if mode.IsDir() { // folder
		// walk through every file in the folder
		filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			// generate tar header
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// must provide real name
			// (see https://golang.org/src/archive/tar/common.go?#L626)
			header.Name = filepath.ToSlash(file)

			// write header
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			// if not a dir, write file content
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err := io.Copy(tw, data); err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		return fmt.Errorf("error: file type not supported")
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}

func (e *Exporter) CreateTmpDirs() {
	dirName, err := e.MkTempDir("tmp")
	if err != nil {
		log.Fatalf("Failed to create tempDir: %v", err)
	}
	e.TmpDir = dirName
	e.TmpAttachmentDir = filepath.Join(e.TmpDir, "attachments")
	e.TmpRepositoryDir = filepath.Join(e.TmpDir, "repositories")

	err = e.Mkdirp(e.TmpDir)
	if err != nil {
		log.Fatalf("Failed to create base migration directory: %v", err)
	}

	err = e.Mkdirp(e.TmpAttachmentDir)
	if err != nil {
		log.Fatalf("Failed to create attachments directory: %v", err)
	}

	err = e.Mkdirp(e.TmpRepositoryDir)
	if err != nil {
		log.Fatalf("Failed to create repositories directory: %v", err)
	}
}

func (e *Exporter) DownloadFile(path string) error {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	dir, file := filepath.Split(path)
	downloadsDir := filepath.Join(e.TmpAttachmentDir, dir)
	// Create dirs for attachments, i.e., `{cwd}/migration/attachments/uploads/1e690825814e23a29bd8810f567829e7`.
	err := e.Mkdirp(downloadsDir)
	if err != nil {
		return err
	}

	resp, err := client.Get(fmt.Sprintf("%s/%s", e.CurrentProject.WebURL, path))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath.Join(downloadsDir, file), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	contentType, err := e.GetContentType(filepath.Join(downloadsDir, file))
	if err != nil {
		return err
	}

	fmt.Println("contentType", contentType)

	_, err = io.Copy(f, resp.Body)
	return err
}

func (e *Exporter) Export() error {
	var wg sync.WaitGroup
	wg.Add(len(e.Exports))

	for _, export := range e.Exports {
		go func(export []string) {
			namespace := export[0]
			projectName := export[1]
			projects, _, err := e.Groups.GetGroupProject(namespace, projectName)
			if err != nil {
				log.Fatalf("Failed to get projects: %v", err)
			}

			if len(projects) > 0 {
				for _, project := range projects {
					_, err := e.Repositories.Clone(project)
					if err != nil {
						if err != nil {
							log.Fatalf("Failed to clone repository `%s`: %v", project.Name, err)
						}
					}

					e.CurrentProject = project

					//e.Projects.Export(project.ID)
					e.Branches.Export()
					e.CommitComments.Export()
					e.Issues.Export()
					e.Groups.Export(namespace)
					e.Labels.Export()
//					e.MergeRequests.Export()
					e.Milestones.Export()
					e.Tags.Export()
					e.Users.Export()
					e.Repositories.Export(project)
				}

			}
			wg.Done()
		}(export)
	}
	wg.Wait()

	e.GetState()

	return nil
}

func (e *Exporter) GetContentType(filename string) (*string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buffer)

	return &contentType, nil
}

func (e *Exporter) GetState() {
	if len(e.State.Branches) > 0 {
		e.Branches.WriteFile()
	}

	if len(e.State.CommitComments) > 0 {
		e.CommitComments.WriteFile()
	}

	if len(e.State.Issues) > 0 {
		e.Issues.WriteFile()
	}

	if len(e.State.Groups) > 0 {
		e.Groups.WriteFile()
	}

	if len(e.State.MergeRequests) > 0 {
		e.MergeRequests.WriteFile()
	}

	if len(e.State.Milestones) > 0 {
		e.Milestones.WriteFile()
	}

	if len(e.State.Repositories) > 0 {
		e.Repositories.WriteFile()
	}

	if len(e.State.Tags) > 0 {
		e.Tags.WriteFile()
	}

	if len(e.State.Users) > 0 {
		e.Users.WriteFile()
	}
}

func (e *Exporter) Mkdirp(path string) error {
	return os.MkdirAll(path, 0755)
}

func (e *Exporter) MkTempDir(tmpDir string) (string, error) {
	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return "", err
		}
	}
	return ioutil.TempDir(tmpDir, "gitlab-exporter")
}

func (e *Exporter) NewRequest(api string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", e.BaseAPI, api)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Set("Accept", string(e.Token))
	headers.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//	defer resp.Body.Close()
	return resp, nil
}

func (e *Exporter) WriteFile(filename string, b []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(b); err != nil {
		f.Close() // ignore captured `err` error from the call to `.Write()`; Write error takes precedence
		log.Fatal(err)
	}
	return f.Close()
}

func (e *Exporter) WriteJsonFile(filename string, state interface{}) error {
	b, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	// Write the file into the `e.BaseDir` migration/archive directory.
	return e.WriteFile(filepath.Join(e.TmpDir, filename), b)
}

func (e *Exporter) WriteStaticFiles() error {
	// static/urls.json
	// static/schema.json -> {version => "1.2.0"}
	entries, err := ioutil.ReadDir(e.StaticDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		f, err := os.Open(filepath.Join(e.StaticDir, entry.Name()))
		if err != nil {
			log.Fatalf("Failed to open file entry `%s`: %v", entry.Name, err)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("Failed to read file entry `%s`: %v", entry.Name, err)
		}
		err = e.WriteFile(filepath.Join(e.TmpDir, entry.Name()), b)
		if err != nil {
			log.Fatalf("Failed to write file entry `%s`: %v", entry.Name, err)
		}
	}

	return nil
}
