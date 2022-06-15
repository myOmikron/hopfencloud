
function downloadFile(fileID) {
    let url = new URL(window.location.toString());
    url.pathname = "/download";
    url.searchParams.set("id", fileID);
    let link = document.createElement("a");
    link.href = url.toString();
    link.click();
}

function changeDirectory(directoryID) {
    let url = new URL(window.location.toString());
    if (directoryID !== undefined) {
        url.searchParams.set("dir", directoryID);
    } else {
        url.searchParams.delete("dir");
    }
    window.location.assign(url);
}

function sendFile(file, progressFunc, loadFunc) {
    let data = new FormData();
    data.append('file', file);

    let request = new XMLHttpRequest();
    request.open('POST', '/upload');
    if (progressFunc !== undefined) {
        request.upload.onprogress = progressFunc
    }
    if (progressFunc !== undefined) {
        request.upload.onloadend = loadFunc;
    }

    request.send(data);
}

function constructFileEntry(file) {
    let li = document.createElement("li");
    li.innerText = file.name;
    return li;
}

window.onload = function() {
    const dropzone = document.getElementById("dropzone");
    const dropzone_fl = document.getElementById("dropzone-fl");
    const dropzone_fs = document.getElementById("dropzone-fs");
    const dropzone_upload = document.getElementById("dropzone-upload");

    let upload_mutex = false;

    dropzone.ondragover = dropzone.ondragenter = (event) => {
        event.stopPropagation();
        event.preventDefault();
        dropzone.classList.add("dropActive");
    };

    dropzone.ondragleave = (event) => {
        event.stopPropagation();
        event.preventDefault();
        dropzone.classList.remove("dropActive");
    };

    dropzone.ondrop = (event) => {
        event.stopPropagation();
        event.preventDefault();
        if (upload_mutex) return;

        dropzone_fs.files = event.dataTransfer.files;
        dropzone_fl.replaceChildren();
        for (const file of event.dataTransfer.files) {
            dropzone_fl.appendChild(constructFileEntry(file));
        }
    };

    dropzone.onclick = () => {
        if (upload_mutex) return;
        dropzone_fs.click();
    };

    dropzone_fs.onchange = () => {
        dropzone_fl.replaceChildren();
        for (const file of dropzone_fs.files) {
            dropzone_fl.appendChild(constructFileEntry(file));
        }
    };

    dropzone_upload.onclick = () => {
        if (upload_mutex) return;
        upload_mutex = true;

        for (let i = 0; i < dropzone_fs.files.length; i++) {
            let child = dropzone_fl.children[i];
            let outside = document.createElement("div");
            let inside = document.createElement("div");
            outside.classList.add("dropzone-progress-border");
            inside.classList.add("dropzone-progress");
            outside.appendChild(inside);
            child.appendChild(outside);

            sendFile(dropzone_fs.files[i], (pe) => {
                if (pe.lengthComputable) {
                    let file_loaded = Math.floor((pe.loaded/pe.total) * 100);
                    inside.style = `width: ${file_loaded}%;`;
                }
            }, () => {
                child.remove();
                if (dropzone_fl.children.length === 0) {
                    window.location.reload();
                }
            });
        }
    };
};
