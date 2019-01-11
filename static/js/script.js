function browsePage(){
        window.location='/repo/';
    }

document.getElementById('fileinputMulti').addEventListener('change', function(){
    for(var i = 0; i<this.files.length; i++){
        var file =  this.files[i];
        // This code is only for demo ...
        //console.group("File "+i);
        //console.log("name : " + file.name);
        //console.log("size : " + file.size);
        //console.log("type : " + file.type);
        //console.log("date : " + file.lastModified);
        //console.groupEnd();

        var formData = new FormData();
        formData.append('file', file);
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/upload", true);
        xhr.send(formData);
    }
}, false);






