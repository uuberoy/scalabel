***Note: This code base is under-going dramatic changes in preparation of our first release slated by the end of June.***

# ScaLabel

[![Build Status](https://travis-ci.com/ucbdrive/scalabel.svg?token=9QKS6inVmkjyhrWUHjqT&branch=master)](https://travis-ci.com/ucbdrive/scalabel)

ScaLabel is a versatile and scalable tool that supports all kinds of annotations needed in a driving database, such as bounding box, semantic instance segmentation, and lane detection.
A screenshot of bounding box annotation interface:
![alt text](/example/bbox_tool.jpg)

## Setup ##
1. Checkout the code
    ```
    git clone git@github.com:ucbdrive/scalabel.git
    cd scalabel
    ```
2. Create a directory <data_dir> to store the server data:
    ```
    mkdir ../data
    ```
3. Launch server. There are two options, either (i) to build 
 with Docker or (ii) to build by yourself.
    1. Build and run a Docker image from the Dockerfile.

        Build by yourself
        ```
        docker build . -t 'scalabel/server'
        ```
        Or
        ```
        docker pull scalabel/server
        ```
        After getting the docker image, you can run the server
        ```
        docker run -it -v `pwd`/../data:/go/data -p 8686:8686 scalabel/server
        ```
    2. Build the server by yourself. 
        1. Install GoLang. Refer to the [instruction page](https://golang.org/doc/install) for details.
        2. Install GoLang dependency  
        ```
        go get gopkg.in/yaml.v2
        ```
        3. Compile the packages 
        ```
        go build -i -o bin/scalabel ./server/go
        ```
        4. Specify basic configurations (e.g. the port to start the server, 
        the data directory to store the image annotations, etc) in your own 
        `config.yml`. Refer to `app/config/default_config.yml` for the default configurations. 
        5. Launch the server by running 
        ```
        ./bin/scalabel --config app/config/default_config.yml
        ```
    
3. Access the server through the specifed port (we use `8686` as the default port
specified in the `config.yml`)
    ```
    http://localhost:8686
    ```

## Task Configurations ##

To start your own annotation project, you need to create the following
configuration files to specify the list of object categories you wish to label, 
the paths to the image files, etc. We will navigate you through the tool in 
the next section.

* Create **categories.yml** with a list of object categories you wish to label. 
Refer to examples in `/example/categories.yml`. 
* Create **image_list.yml** with a list of paths to image files. These images 
should be publicly accessible or stored in the server. Refer to examples in
`/example/image_list.yml`.
* (Optional) Create **attributes.yml** to 
describe more fine-grained and intra-class variation
(e.g. is the person walking or standing, or what is the traffic light color).
Refer to examples in `example/bbox_attributes.yml`.

## Navigating the tool ##
We are now navigating you through the tool to demonstrate how to create a project 
and monitor the project progress in the dashboard page.

### Project Creation ###

Open http://localhost:8686/create in your browser to create a new task. 
You'll see several attributes described as follows in the interface. 

* **Project name**: A customized name given to your project. Upon finishing creating the project, you can find all task assignments in 
`<data_dir>/Assignments/<project_name>/`, completed annotations in  
`<data_dir>/Submissions/<project_name>/`, and log files in 
`<data_dir>/Logs/<project_name>/` later on.
* **Item type**: Specify the data type of your project. Choose one of the three item types supported for now: 
`Image`, `Video` and `Point Cloud`.
* **Label type**: Specify the label type of your project. 
Choose one of the five label types: '2d_bbox', '2d_road', '2d_seg', '2d_lane' and 'image'.
* **Page Title**: The page title of the dashboard page. A default title
will be filled in once you specify the `Label Type`.
* **Item List**: Upload the created `image_list.yml` which contains a list of paths to image files. 
(Input format refers to `example/image_list.yml`.)
* **Categories**: Upload the created `categories.yml` which specifies the label categories (object classes) such as car, person and truck etc. 
(Input format refers to `example/categories.yml`.)
* **Attributes**: (Optional) Upload the created `attributes.yml` which describes more fine-grained and intra-class variation (e.g. is the person walking or standing, or what is the traffic light color). 
(Input format refers to `example/bbox_attributes.yml`)
* **Task size**: Specify the number of images to be labeled for each assignment. Based on the number of total images and the task size, our tool will automatically decide the number of assignments.
* **Vendor ID**: This keeps track of which vendor/annotator to be assigned to. 
Specify `0` if unsure.


### Dashboard ###

After specifying the attributes above, click '**enter**'and 
then '**go to dashboard**'. It will direct you to the dashboard --
http://localhost:8686/dashboard?project_name=<project_name>... From here, you can 
* Access and monitor the list of tasks created. 
* Download a list of task URLs, and send them to your vendor. Vendors can access and work on URLs concurrently. 
* Download all results of annotation anytime.


### Workflow ###

Here is the labeling tool workflow diagram:
   
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
