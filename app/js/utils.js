// Record timestamp in milliseconds, action, and target image index
function addEvent(action, index, position) {
    if (!assignment.events) {
        assignment.events = [];
    }
    var event = {
        "timestamp": Math.round(new Date() / 1000),
        "action": action,
        "targetIndex": index.toString(),
        "position": position // only applicable to certain actions
    };
    assignment.events.push(event);
}

// Preload images using browser caching
function preload(imageArray, index) {
    index = index || 0;
    if (imageArray && imageArray.length > index) {
        preloaded_images[index] = new Image();
        preloaded_images[index].onload = function () {
            // addEvent("image loaded", index);
            if (index === 0) {
                // display when the first image is loaded
                bboxLabeling = new BBoxLabeling({
                    url: preloaded_images[current_index].src
                });
                bboxLabeling.replay();
                num_display = num_display + 1;
                //addEvent("display", index);
            }
            preload(imageArray, index + 1);
        };
        preloaded_images[index].onerror = function () {
            addEvent("image fails to load", index);
            preload(imageArray, index + 1);
        };
        preloaded_images[index].src = imageArray[index].url;
    } else {
        //console.log("finish preloading all images.");
    }
}

function updateProgressBar() {
    var progress = $("#progress");
    progress.html(" " + (current_index + 1).toString() + "/" +
        assignment.taskSize.toString())
}


function updateCategorySelect() {
    var category = assignment.category;
    var category_select = $("select#category_select");

    for (var i = 0; i < category.length; i++) {
        if (category[i]) {
            category_select.append("<option>" +
                category[i] + "</option>");
        }
    }
    $("select#category_select").val(assignment.category[0]);
}

// Update global image list
function saveLabels() {
    bboxLabeling.submitLabels();
    image_list[current_index].labels = bboxLabeling.output_labels;
    image_list[current_index].tags = bboxLabeling.output_tags;
}

function submitAssignment() {
    var x = new XMLHttpRequest();
    x.onreadystatechange = function () {
        if (x.readyState === 4) {
            //console.log(x.response)
        }
    };
    assignment.images = image_list;
    assignment.numLabeledImages = current_index + 1;
    assignment.userAgent = navigator.userAgent;

    x.open("POST", "/postSubmission");
    x.send(JSON.stringify(assignment));
}

function submitLog() {
    var x = new XMLHttpRequest();
    x.onreadystatechange = function () {
        if (x.readyState === 4) {
            //console.log(x.response)
        }
    };
    assignment.images = image_list;
    assignment.numLabeledImages = current_index + 1;
    assignment.userAgent = navigator.userAgent;

    x.open("POST", "/postLog");
    x.send(JSON.stringify(assignment));
}

function loadAssignment() {
    var x = new XMLHttpRequest();
    x.onreadystatechange = function () {
        if (x.readyState === 4) {
            //console.log(x.response);
            assignment = JSON.parse(x.response);
            image_list = assignment.images;
            current_index = 0;
            addEvent("start labeling", current_index);
            assignment.startTime = Math.round(new Date() / 1000);

            // preload images
            preload(image_list);

            getIPAddress();
            // update toolbar and progress bar
            updateCategorySelect();
            updateProgressBar();
        }
    };

    // get params from url path
    var searchParams = new URLSearchParams(window.location.search);
    var task_id = searchParams.get('task_id');
    var project_name = searchParams.get('project_name');

    var request = JSON.stringify({
        "assignmentId": task_id,
        "projectName": project_name
    });

    x.open("POST", "/requestSubmission");
    x.send(request);
}

function getIPAddress() {
    $.getJSON('//ipinfo.io/json', function (data) {
        assignment.ipAddress = data;
    });
}

function goToImage(index) {

    saveLabels();
    // auto save log every twenty five images displayed
    if (num_display % 25 === 0 && num_display !== 0) {
        submitLog();
        addEvent("save log", index);
    }
    // auto save submission for the current session.
    submitAssignment();

    if (index === -1) {
        alert("This is the first image.");
    } else if (index === image_list.length) {
        addEvent("submit", index);
        alert("Good Job! You've completed this assignment.");
    } else {
        current_index = index;
        num_display = num_display + 1;
        addEvent("save", index);
        if (index === image_list.length - 1) {
            $('#save_btn').text("Submit");
            $('#save_btn').removeClass("btn-primary").addClass("btn-success");
        }
        if (index === image_list.length - 2) {
            $('#save_btn').removeClass("btn-success").addClass("btn-primary");
            $('#save_btn').text("Save");
        }
        addEvent("display", index);
        updateProgressBar();
        bboxLabeling.updateImage(preloaded_images[index].src);
        bboxLabeling.replay();
    }
}
