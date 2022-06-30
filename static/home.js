let imageInput = document.querySelector("#image-upload");

document.querySelector("#javascriptContent").style.display='block';
document.querySelector("#nonJavascriptContent").style.display='none';

document.addEventListener('click', function() {
   canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height);
   let image_data_url = canvas.toDataURL('image/jpeg');

   // data url of the image
   console.log(image_data_url);
   appendFileAndSubmit(image_data_url);
});

imageInput.addEventListener('change', function() {
   document.querySelector("#image-form").submit();
});

let video = document.querySelector("#video");
let canvas = document.querySelector("#canvas");

navigator.mediaDevices.getUserMedia({
   video: true,
   audio: false
}).then(function(stream) {
   video.srcObject = stream;
}).catch(function(err) {
   /* handle the error */
   console.log("something went wrong while getting access to the camera: " + err);
});

function appendFileAndSubmit(ImageURL){
   // Get the form
   var form = document.getElementById("desktopForm");

   // Split the base64 string in data and contentType
   var block = ImageURL.split(";");
   // Get the content type
   var contentType = block[0].split(":")[1];
   // get the real base64 content of the file
   var realData = block[1].split(",")[1];
   console.log(contentType);
   // Convert to blob
   var blob = b64toBlob(realData, contentType);

   // Create a FormData and append the file
   var fd = new FormData(form);
   fd.append("image-file", blob);

   // Submit Form and upload file
   $.ajax({
       url:"/upload",
       data: fd,
       type:"POST",
       contentType:false,
       processData:false,
       cache:false,
       dataType:"json", // Change this according to your response from the server.
       error:function(err){
           console.error(err);
       },
       success:function(data){
           console.log(data);
       },
       complete:function(){
           console.log("Request finished.");
       }
   });
}

function b64toBlob(b64Data, contentType, sliceSize) {
   contentType = contentType || '';
   sliceSize = sliceSize || 512;

   var byteCharacters = atob(b64Data);
   var byteArrays = [];

   for (var offset = 0; offset < byteCharacters.length; offset += sliceSize) {
       var slice = byteCharacters.slice(offset, offset + sliceSize);

       var byteNumbers = new Array(slice.length);
       for (var i = 0; i < slice.length; i++) {
           byteNumbers[i] = slice.charCodeAt(i);
       }

       var byteArray = new Uint8Array(byteNumbers);

       byteArrays.push(byteArray);
   }

 var blob = new Blob(byteArrays, {type: contentType});
 return blob;
}
