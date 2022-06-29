let video = document.querySelector("#video");
let canvas = document.querySelector("#canvas");

var ResolutionsToCheck = [
                {width: 160, height:120},
                {width: 320, height:180},
                {width: 320, height:240},
                {width: 640, height:360},
                {width: 640, height:480},
                {width: 768, height:576},
                {width: 1024, height:576},
                {width: 1280, height:720},
                {width: 1280, height:768},
                {width: 1280, height:800},
                {width: 1280, height:900},
                {width: 1280, height:1000},
                {width: 1920, height:1080},
                {width: 1920, height:1200},
                {width: 2560, height:1440},
                {width: 3840, height:2160},
                {width: 4096, height:2160}
            ];

var left = 0;
var right = ResolutionsToCheck.length;
var selectedWidth;
var selectedHeight;
var mid;

function FindMaximum_WidthHeight_ForCamera()
{
    console.log("left:right = ", left, ":", right);
    if(left >= right)
    {
        console.log("Selected Height:Width = ", selectedWidth, ":", selectedHeight);
        return;
    }

    mid = Math.floor((left + right) / 2);

    var temporaryConstraints = {
        "audio": true,
        "video": {
            "mandatory": {
            "minWidth": ResolutionsToCheck[mid].width,
            "minHeight": ResolutionsToCheck[mid].height,
            "maxWidth": ResolutionsToCheck[mid].width,
            "maxHeight": ResolutionsToCheck[mid].height
            },
        "optional": []
        }
    }

    navigator.mediaDevices.getUserMedia(temporaryConstraints).then(checkSuccess).catch(checkError);
}

function checkSuccess(stream)
{
    console.log("Success for --> " , mid , " ", ResolutionsToCheck[mid]);
    selectedWidth = ResolutionsToCheck[mid].width;
    selectedHeight = ResolutionsToCheck[mid].height;

    left = mid+1;

    for (let track of stream.getTracks()) {
        track.stop()
    }

    FindMaximum_WidthHeight_ForCamera();
}


function checkError(error)
{
    console.log("Failed for --> " + mid , " ", ResolutionsToCheck[mid],  " ", error);
    right = mid-1;

    FindMaximum_WidthHeight_ForCamera();
}

FindMaximum_WidthHeight_ForCamera();

navigator.mediaDevices.getUserMedia({
      video: { width: selectedWidth, height: selectedHeight },
      audio: false
}).then(function(stream) {
  /* use the stream */
   video.srcObject = stream;
}).catch(function(err) {
  /* handle the error */
   console.log("error: " + err);
});



document.addEventListener('click', function() {
   console.log("clicked");
   canvas.width = selectedWidth;
   canvas.height = selectedHeight;
   canvas.getContext('2d').drawImage(video, 0, 0, selectedWidth, selectedHeight);
   let image_data_url = canvas.toDataURL('image/jpeg');

   // data url of the image
   console.log(image_data_url);
});
