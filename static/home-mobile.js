let imageInput = document.querySelector("#image-upload");

document.querySelector("#javascriptContent").style.display='block';
document.querySelector("#nonJavascriptContent").style.display='none';

document.addEventListener('click', function() {
   imageInput.click();
});

imageInput.addEventListener('change', function() {
   document.querySelector("#image-form").submit();
});
