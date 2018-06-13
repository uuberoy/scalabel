***Note: This code base is under-going dramatic changes in preparation of our first release slated by the end of June.***

# Scalable Annotation Tooling

[![Build Status](https://travis-ci.com/ucbdrive/sat.svg?token=9QKS6inVmkjyhrWUHjqT)](https://travis-ci.com/ucbdrive/sat)

SAT is a versatile and scalable tool that supports all kinds of annotations needed in a driving database, such as bounding box, semantic instance segmentation, and lane detection.  
A screenshot of bounding box annotation interface:
![alt text](/example/bbox_tool.jpg)

## Initial Setup ##
* Create a directory to store image annotations: 
```
mkdir ../data
```
* (Optional) Build and run a Docker image from the 
Dockerfile. 
e.x.
```
docker build . -t 'scalabel:demo'
docker run -it -p 8686:8686 scalabel:demo
```

## Images and Labels ##
* Create **label.txt** with a list of object categories you wish to label. 
Refer to examples in /example/bbox_label.txt and /example/drivable_label.txt. 
* Create **image_list.json** with a list of paths to image files. These images 
should be publicly accessible or stored in the server. Refer to example in
/example/image_list.json

## Running Server ##
Compiles the packages named by the import paths, along with their dependencies, but it does not install the results.:  
```
go build -i -o bin/sat ./server/go
```
Specify a port to start the server and a directory path to store 
image annotations. Use the following command to run the server:
```
./bin/sat --config config.yml
```

## Navigating the tool ##
We'll demonstrate how to navigate our tool. Here, we default to localhost and 
listen to port 8686.

1. Go to http://localhost:8686/create to create a new task. You need to 
upload label.txt and image_list.json you just created, and specify other 
task configurations such as project name, task size, label type, and vendor ID. 
* **Project name**: Find all task assignments in 
<data_dir>/Assignments/<project_name>/, completed annotations in  
<data_dir>/Submissions/<project_name>/, and log files in 
<data_dir>/Logs/<project_name>/
* **Task size**: Specify number of images to be labeled for each assignment.
* **Label type**: Choose one of the five label types: '2d_bbox', '2d_road', '2d_seg', '2d_lane','image'.
* **Vendor ID**: This keeps track of which vendor/annotator to be assigned to. 
Specify '0' if unsure.
* **Image List**: Specify list of paths to image files. (Input format refers to image_list.json in *example folder*.)
* **Label**: Specify label categories (object classes) such as car, person and truck etc. (Input format refers to ~_label.txt in *example folder*.)
* **Attributes**: Specify attributes which describe more fine-grained and intra-class variation(e.g. is the person walking or standing, or what is the traffic light color). This feature is optional. 
(Input format refers to ~_attributes.json in *example folder*.)

2. Click '**enter**', then '**go to dashboard**' will direct you to the dashboard --
http://localhost:8686/dashboard?project_name=<project_name>... From here, you can 
* Access and monitor the list of tasks created. 
* Download a list of task URLs, and send them to your vendor. Vendors can access and work on URLs concurrently. 
* Download all results of annotation anytime.

Here is a labeling tool workflow diagram:
![alt text](/example/workflow.jpg) 

## Directories ##
* **app folder**: front-end scripts for control panel, bounding box, image level, drivable area, lane markings, and segmentation.
* **example folder**: example tasks.
* **main.go**: back-end go-lang server script.

## How to Use Region Annotation ##
Region Annonation consists of annotations of drivable area, lane markings and image segmentation, corresponding to label types of "2d_road", "2d_lane", and "2d_seg".
* The tool box is located on the left side of the screen. 
* To draw a label, click along the edge of an object to draw a mask, then go back to the first vertex (where you started) to form a closed path.
* Press **ESC** to remove an unfinished object.
* To change the category of an object, **double click** on it and reselect from the toolbox.
* To remove an object, click **delete** (or press **delete** on keyboard).
* To change the shape of an object, drag its vertices. When hovering the mouse on these points, they will become green and bigger.
* To delete a vertex, hover the mouse on it and press **delete**.
* To add a vertex, drag the corresponding midpoint (in pink). The midpoint will become a new vertex after dragging, and two new midpoints will be generated. The corresponding midpoint is in orange when you hover your mouse on it.


A screenshot of vertexs and midpoints:


![alt text](/example/vertex.png)![alt text](/example/midpoint.png)


* To zoom in/out of an image, press **PageUp/PageDown** on keyboard, or use a local magnifier by choosing **Magnify** in the toolbox, or click **+/-** on top of the webpage.  
Bezier Curve:
* To add bezier curve: Hover your mouse on a midpoint, then press "**B**" or "**b**" on the keyboard. The midpoint will split into two points on the corresponding edge, and you can drag these two points and the other two endpoints to change the shape of the Bezier Curve.
* To delete bezier curve: Hover your mouse on a control point of the Bezier Curve, then press **delete** on the keyboard. It will become a straight edge.


## Additional Options for Segmentation Annotation ##
### Share Border 
1. To draw the coincide part of two objects, press "**S**" or "**s**" before or during drawing;
2. Vertices of other objects will become visible when you move your mouse onto them. Meanwhile, the toolbox will become blue. It will remain blue or green unless you exit the "quick draw" mode. (By pressing "S/s" again).
3. Select the first and the last vertices on the border by clicking on exiting vertices. After the first click, the object you selected will change color, and the toolbox will turn green. After the second click, a polyline will be formed, and the toolbox will change from green to blue. 
4. Use **PageLeft/PageRight** on your keyboard to choose the border, since there are two directions: clockwise or counter-clockwise.
5. The toolbox will remain green before you click a pair of points on the same object. And when you click on a new object, it will remain green, since no coincide border is formed, so that you can continue finding its pair.
6. Click the background image, the toolbox will be blue. So in conclusion, color green means you have finished the first click and should look for the second point; color blue means you should look for the first point.
7. The same as ordinary draw mode, press "delete" to delete the latest vertex or last added border, and press "esc" to delete the whole object. When a closed path is formed, your annotation is finished.
8. You can use it to click anywhere like the ordinary draw mode, it only provides function to draw coincide part quickly by clicking **twice**.
### Link Segments  
To link separated parts of an object into one with the same object ID
1. Select one part of the object (by **double click**), and click **Link** button, then single click the others. Click **finish link** when finishes.
2. You can also click **Link** first, then single click every object to be linked. Click **finish link** when finishes.

## Contribute
Before committing code, please run
```
sh scripts/setup.sh
```
