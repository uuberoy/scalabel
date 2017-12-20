# Scalable Anntoation Tooling

## Overview ##
SAT is a versatile and scalable tool that supports all kinds of annotations needed in a driving database, such as bounding box, semantic instance segmentation, and lane detection. Here are our labeling tool workflow diagram and a screenshot of bounding box annotation interface.
![alt text](/example/workflow.png) | ![alt text](/example/bbox_tool.png)

## Initial Setup ##
* Create a directory to store image annotations: 
mkdir ../data

## Images and Labels ##
* Create label.txt with a list of object categories you wish to label. 
Refer to examples in /example/bbox_label.txt and /examples/drivable_label.txt. 
* Create image_list.json with a list of paths to image files. These images 
should be publicly accessible or stored in the server. Refer to example in
/example/image_list.json

## Running Server ##
Specify a port to start the server and a directory path to store 
image annotations. Use the following command to run the server:
go run main.go --port 8686 --data_dir "../data"

## Navigating the tool ##
We'll demonstrate how to navigate our tool. Here, we default to localhost and 
listen to port 8686.

1. Go to http://localhost:8686/create to create a new task. You need to 
upload label.txt and image_list.json you just created, and specify other 
task configurations such as project name, task size, label type, and vendor ID. 
* Project name: You can find all task assignments in 
<data_dir>/Assignments/<project name>/, completed annotations in  
<data_dir>/Submissions/<project name>/, and log files in 
<data_dir>/Logs/<project name>/
* Task size: You can specify number of images to be labeled for each assignment.
* Label type: We provide three different label types: '2d_bbox', '2d_poly', 'image'.
* Vendor ID: This keeps track of which vendor/annotator to be assigned to. 
You can specify '0' if unsure.

2. Once you click 'enter', 'go to dashboard' will direct you to 
http://localhost:8686/dashboard?project_name=<project 
name>... 
* From here, you can access and monitor list of tasks 
created. 
* You can download a list of task URLs, and send them to your 
vendor. Vendors can concurrently access and work on one task URL. 
* You can download all result of annotation at anytime.

## Directories ##
* app folder contains front-end scripts for control panel, bounding box, 
image level, and drivable annotations.
* example folder contains task examples.
* main.go is the back-end go-lang server script.

