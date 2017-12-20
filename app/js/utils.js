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
                if(type == "bbox"){
                    bboxLabeling = new BBoxLabeling({
                        url: preloaded_images[current_index].src
                    });
                    bboxLabeling.replay();
                } else {
                    polyLabeling = new PolyLabeling({
                        url: preloaded_images[current_index].src
                    });
                    polyLabeling.updateImage(preloaded_images[current_index].src);
                }
                num_display = num_display + 1;
            }
            preload(imageArray, index + 1);
        };
        preloaded_images[index].onerror = function () {
            addEvent("image fails to load", index);
            preload(imageArray, index + 1);
        };
        preloaded_images[index].src = imageArray[index].url;
    } else {
        $("#prev_btn").attr("disabled", false);
        $("#next_btn").attr("disabled", false);
    }
}

function updateProgressBar() {
    var progress = $("#progress");
    progress.html(" " + (current_index + 1).toString() + "/" +
        assignment.taskSize.toString())
}


function updateCategorySelect() {
    if (type == "poly"){
        updateCategory();
    } else {
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
}

// Update global image list
function saveLabels() {
    if(type == "bbox") {
        bboxLabeling.submitLabels();
        image_list[current_index].labels = bboxLabeling.output_labels;
        image_list[current_index].tags = bboxLabeling.output_tags;
    } else {
        polyLabeling.submitLabels();
        image_list[current_index].labels = polyLabeling.output_labels;
    }
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
            if (type == "poly") {
                for (var idx in image_list) {
                    var labels = image_list[idx].labels;
                    for (var key in labels) {
                        if (labels.hasOwnProperty(key)) {
                            var label = labels[key];
                            num_poly = num_poly + 1;
                        }
                    }
                }
                $("#poly_count").text(num_poly);
            }
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
        if(type == "bbox") {
            bboxLabeling.updateImage(preloaded_images[index].src);
            bboxLabeling.replay();
        } else {
            polyLabeling.clearAll();
            polyLabeling.updateImage(preloaded_images[index].src);
        }
    }
}
