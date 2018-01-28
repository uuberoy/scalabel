package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Task struct {
	AssignmentID       string        `json:"assignmentId"`
	ProjectName        string        `json:"projectName"`
	WorkerID           string        `json:"workerId"`
	Category           []string      `json:"category"`
	LabelType          string        `json:"labelType"`
	TaskSize           int           `json:"taskSize"`
	Images             []ImageObject `json:"images"`
	SubmitTime         int64         `json:"submitTime"`
	NumSubmissions     int           `json:"numSubmissions"`
	NumLabeledImages   int           `json:"numLabeledImages"`
	NumDisplayedImages int           `json:"numDisplayedImages"`
	StartTime          int64         `json:"startTime"`
	Events             []Event       `json:"events"`
	VendorID           string        `json:"vendorId"`
	IPAddress          interface{}   `json:"ipAddress"`
	UserAgent          string        `json:"userAgent"`
}

type Result struct {
	Images []ImageObject `json:"images"`
}

type TaskInfo struct {
	AssignmentID     string `json:"assignmentId"`
	ProjectName      string `json:"projectName"`
	WorkerID         string `json:"workerId"`
	LabelType        string `json:"labelType"`
	TaskSize         int    `json:"taskSize"`
	SubmitTime       int64  `json:"submitTime"`
	NumSubmissions   int    `json:"numSubmissions"`
	NumLabeledImages int    `json:"numLabeledImages"`
	StartTime        int64  `json:"startTime"`
}

type Event struct {
	Timestamp   int64       `json:"timestamp"`
	Action      string      `json:"action"`
	TargetIndex string      `json:"targetIndex"`
	Position    interface{} `json:"position"`
}

type ImageObject struct {
	Url         string   `json:"url"`
	GroundTruth string   `json:"groundTruth"`
	Labels      []Label  `json:"labels"`
	Tags        []string `json:"tags"`
}

type Label struct {
	Id        string      `json:"id"`
	Category  string      `json:"category"`
	Attribute interface{} `json:"attribute"`
	Position  interface{} `json:"position"`
}

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	port     = flag.String("port", "", "")
	data_dir = flag.String("data_dir", "", "")
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)

	flag.StringVar(port, "s", "8686", "")
	flag.StringVar(data_dir, "d", "../data", "")
}

var HTML []byte
var mux *http.ServeMux

func (assignment *Task) GetAssignmentPath() string {
	filename := assignment.AssignmentID
	dir := path.Join(*data_dir, "Assignments", assignment.ProjectName)
	os.MkdirAll(dir, 0777)
	return path.Join(dir, filename+".json")
}

func (assignment *Task) GetSubmissionPath() string {
	start_time := formatTime(assignment.StartTime)
	dir := path.Join(*data_dir, "Submissions", assignment.ProjectName,
		assignment.AssignmentID)
	os.MkdirAll(dir, 0777)
	return path.Join(dir, start_time+".json")
}

func (assignment *Task) GetLatestSubmissionPath() string {
	dir := path.Join(*data_dir, "Submissions", assignment.ProjectName,
		assignment.AssignmentID)
	os.MkdirAll(dir, 0777)
	return path.Join(dir, "latest.json")
}

func (assignment *Task) GetLogPath() string {
	submit_time := formatTime(assignment.SubmitTime)
	dir := path.Join(*data_dir, "Log", assignment.ProjectName,
		assignment.AssignmentID)
	os.MkdirAll(dir, 0777)
	return path.Join(dir, submit_time+".json")
}

func recordTimestamp() int64 {
	// record timestamp in seconds
	return time.Now().Unix()
}

func formatTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02_03-04-05")
}

func formatID(id int) string {
	str := strconv.Itoa(id)
	str_len := utf8.RuneCountInString(str)
	for i := 0; i < (4 - str_len); i += 1 {
		str = "0" + str
	}
	return str
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	flag.Parse()
	// Mux for static files
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./app")))

	// routes
	http.HandleFunc("/", parse(indexHandler))
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/2d_bbox_labeling", bboxLabelingHandler)
	http.HandleFunc("/2d_road_labeling", roadLabelingHandler)
	http.HandleFunc("/2d_seg_labeling", segLabelingHandler)
	http.HandleFunc("/2d_lane_labeling", laneLabelingHandler)
	http.HandleFunc("/image_labeling", imageLabelingHandler)

	http.HandleFunc("/result", readResultHandler)
	http.HandleFunc("/fullResult", readFullResultHandler)

	http.HandleFunc("/postAssignment", postAssignmentHandler)
	http.HandleFunc("/postSubmission", postSubmissionHandler)
	http.HandleFunc("/postLog", postLogHandler)
	http.HandleFunc("/requestAssignment", requestAssignmentHandler)
	http.HandleFunc("/requestSubmission", requestSubmissionHandler)
	http.HandleFunc("/requestInfo", requestInfoHandler)

	http.ListenAndServe(":"+*port, nil)
}

func parse(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if strings.ContainsRune(r.URL.Path, '.') {
			mux.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(HTML)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/control/create.html")
	w.Write(HTML)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/control/monitor.html")
	w.Write(HTML)
}

func bboxLabelingHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/annotation/box.html")
	w.Write(HTML)
}

func roadLabelingHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/annotation/road.html")
	w.Write(HTML)
}

func segLabelingHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/annotation/seg.html")
	w.Write(HTML)
}

func laneLabelingHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/annotation/lane.html")
	w.Write(HTML)
}

func imageLabelingHandler(w http.ResponseWriter, r *http.Request) {
	HTML, _ = ioutil.ReadFile("./app/annotation/image.html")
	w.Write(HTML)
}

func postAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	var task = Task{}
	// Process image list file
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image_list")
	defer file.Close()
	json.NewDecoder(file).Decode(&task)

	// Process label categories file
	label_file, _, err := r.FormFile("label")
	var labels []string
	scanner := bufio.NewScanner(label_file)
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	task_size, err := strconv.Atoi(r.FormValue("task_size"))
	task.ProjectName = r.FormValue("project_name")

	size := len(task.Images)
	assignment_id := 0
	for i := 0; i < size; i += task_size {

		// Initialize new assignment
		assignment := Task{
			ProjectName:      r.FormValue("project_name"),
			LabelType:        r.FormValue("label_type"),
			Category:         labels,
			VendorID:         r.FormValue("vendor_id"),
			AssignmentID:     formatID(assignment_id),
			WorkerID:         strconv.Itoa(assignment_id),
			NumLabeledImages: 0,
			NumSubmissions:   0,
			StartTime:        recordTimestamp(),
			Images:           task.Images[i:Min(i+task_size, size)],
			TaskSize:         task_size}

		assignment_id = assignment_id + 1

		// Save assignment to data folder
		assignment_path := assignment.GetAssignmentPath()

		assignmentJson, _ := json.MarshalIndent(assignment, "", "  ")
		err = ioutil.WriteFile(assignment_path, assignmentJson, 0644)

		if err != nil {
			Error.Println("Failed to save assignment file of",
				assignment.ProjectName, assignment.AssignmentID)
		} else {
			Info.Println("Saving assignment file of",
				assignment.ProjectName, assignment.AssignmentID)
		}
	}

	Info.Println("Created", assignment_id, "new assignments")

	w.Write([]byte(strconv.Itoa(assignment_id)))
}

func postSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Failed to read submission request body")
	}
	assignment := Task{}
	err = json.Unmarshal(body, &assignment)
	if err != nil {
		Error.Println("Failed to parse submission JSON")
	}

	if assignment.NumLabeledImages == assignment.TaskSize {
		assignment.NumSubmissions = assignment.NumSubmissions + 1
		Info.Println("Complete submission of",
			assignment.ProjectName, assignment.AssignmentID)
	}
	assignment.SubmitTime = recordTimestamp()

	submission_path := assignment.GetSubmissionPath()
	taskJson, _ := json.MarshalIndent(assignment, "", "  ")
	err = ioutil.WriteFile(submission_path, taskJson, 0644)
	if err != nil {
		Error.Println("Failed to save submission file of",
			assignment.ProjectName, assignment.AssignmentID)
	}

	latest_submission_path := assignment.GetLatestSubmissionPath()
	latest_taskJson, _ := json.MarshalIndent(assignment, "", "  ")
	err = ioutil.WriteFile(latest_submission_path, latest_taskJson, 0644)
	if err != nil {
		Error.Println("Failed to save latest submission file of",
			assignment.ProjectName, assignment.AssignmentID)
	}
	// Debug
	w.Write(taskJson)

}

func postLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Failed to read log request body")
	}
	assignment := Task{}
	err = json.Unmarshal(body, &assignment)
	if err != nil {
		Error.Println("Failed to parse log JSON")
	}

	if assignment.NumLabeledImages == assignment.TaskSize {
		assignment.NumSubmissions = assignment.NumSubmissions + 1
	}

	assignment.SubmitTime = recordTimestamp()
	// Save to Log every five images displayed
	log_path := assignment.GetLogPath()
	taskJson, _ := json.MarshalIndent(assignment, "", "  ")
	err = ioutil.WriteFile(log_path, taskJson, 0644)
	if err != nil {
		Error.Println("Failed to save log file of",
			assignment.ProjectName, assignment.AssignmentID)
	} else {
		Info.Println("Saving log of",
			assignment.ProjectName, assignment.AssignmentID)
	}

	w.Write(taskJson)
}

func requestAssignmentHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Failed to read assignment request body")
	}
	task := Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		Error.Println("Failed to parse assignment request JSON")
	}
	request_path := task.GetAssignmentPath()

	requestJson, err := ioutil.ReadFile(request_path)
	if err != nil {
		Error.Println("Failed to read assignment file of",
			task.ProjectName, task.AssignmentID)
	} else {
		Info.Println("Finished reading assignment file of",
			task.ProjectName, task.AssignmentID)
	}
	w.Write(requestJson)

}

func requestSubmissionHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Failed to read submission request body")
	}
	request := Task{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		Error.Println("Failed to parse submission request JSON")
	}
	request_path := request.GetLatestSubmissionPath()
	assignment_path := request.GetAssignmentPath()

	var existing_path string
	if Exists(request_path) {
		existing_path = request_path
	} else if Exists(assignment_path) {
		existing_path = assignment_path
	} else {
		Error.Println("Can not find",
			request.ProjectName, request.AssignmentID)
		http.NotFound(w, r)
		return
	}

	requestJson, err := ioutil.ReadFile(existing_path)
	if err != nil {
		Error.Println("Failed to read submission file of",
			request.ProjectName, request.AssignmentID)
	} else {
		Info.Println("Loading assignment from latest submission of",
			request.ProjectName, request.AssignmentID)
	}
	w.Write(requestJson)

}

func requestInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Failed to read submission request body")
	}
	request := Task{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		Error.Println("Failed to parse submission request JSON")
	}
	request_path := request.GetLatestSubmissionPath()
	assignment_path := request.GetAssignmentPath()

	var existing_path string
	if Exists(request_path) {
		existing_path = request_path
	} else if Exists(assignment_path) {
		existing_path = assignment_path
	} else {
		Error.Println("Can not find", assignment_path,
			request.ProjectName, request.AssignmentID)
		http.NotFound(w, r)
		return
	}

	requestJson, err := ioutil.ReadFile(existing_path)
	if err != nil {
		Error.Println("Failed to read submission file of",
			request.ProjectName, request.AssignmentID)
	} else {
		Info.Println("Loading task info of",
			request.ProjectName, request.AssignmentID)
	}

	info := TaskInfo{}
	json.Unmarshal(requestJson, &info)

	infoJson, _ := json.MarshalIndent(info, "", "  ")
	w.Write(infoJson)

}

func readResultHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filename := queryValues.Get("task_id")
	project_name := queryValues.Get("project_name")

	HTML = getResult(filename, project_name)
	w.Write(HTML)
}

func readFullResultHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	project_name := queryValues.Get("project_name")

	HTML = getFullResult(project_name)
	w.Write(HTML)
}

func getResult(assignment_id string, project_name string) []byte {
	submissionPath := path.Join(*data_dir, "Submissions", project_name,
		assignment_id, "latest.json")
	assignment_path := path.Join(*data_dir, "Assignments",
		project_name,
		assignment_id+".json")

	result := Result{}

	var existing_path string
	if Exists(submissionPath) {
		existing_path = submissionPath
	} else if Exists(assignment_path) {
		existing_path = assignment_path
	}

	if len(existing_path) > 0 {
		taskJson, err := ioutil.ReadFile(existing_path)
		if err != nil {
			Error.Println("Failed to read result of",
				project_name, assignment_id)
		} else {
			Info.Println("Reading result of",
				project_name, assignment_id)
		}

		task := Task{}
		json.Unmarshal(taskJson, &task)
		result.Images = task.Images
	}
	resultJson, _ := json.MarshalIndent(result, "", "  ")

	return resultJson
}

func getFullResult(project_name string) []byte {
	result := Result{}
	dir := path.Join(*data_dir, "Submissions", project_name)
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		filename := f.Name()
		submissionPath := path.Join(dir, filename, "latest.json")
		assignment_path := path.Join(*data_dir, "Assignments",
			project_name, filename+".json")

		var existing_path string
		if Exists(submissionPath) {
			existing_path = submissionPath
		} else if Exists(assignment_path) {
			existing_path = assignment_path
		}

		if len(existing_path) > 0 {

			resultJson, err := ioutil.ReadFile(existing_path)
			if err != nil {
				Error.Println("Failed to read result of", project_name)
			} else {
				Info.Println("Reading result of", project_name)
			}

			task := Task{}
			json.Unmarshal(resultJson, &task)
			for i := 0; i < len(task.Images); i += 1 {
				result.Images = append(result.Images, task.Images[i])
			}
		}

	}
	fullResultJson, _ := json.MarshalIndent(result, "", "  ")

	return fullResultJson
}
