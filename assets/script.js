const $ = document.getElementById.bind(document);

bind("image-size", "image-size-label", "Size");
bind("thumb-size", "thumb-size-label", "Size");
bind("quality", "quality-label", "Quality");

$("choose").addEventListener("click", () => {
    choose().then(path => {
        $("directory").value = path;
    })
});

$("start").addEventListener("click", () => {
    let progress = $("progress");
    let imageSize = parseInt(document.getElementById("image-size").value);
    let thumbSize = parseInt(document.getElementById("thumb-size").value);
    let quality = parseInt(document.getElementById("quality").value);
    let directory = document.getElementById("directory").value;

    lock(true);
    progress.hidden = false;
    resize(directory, imageSize, thumbSize, quality).finally(() => {
        lock(false);
        progress.value = 0;
        progress.hidden = true;
    });
});

function bind(input, label, prefix) {
    $(label).innerText = `${prefix}: ` + $(input).value;

    $(input).addEventListener("input", function () {
        $(label).innerText = `${prefix}: ${this.value}`
    });
}

function lock(value) {
    $("image-size").disabled = value;
    $("thumb-size").disabled = value;
    $("quality").disabled = value;
    $("directory").disabled = value;
    $("choose").disabled = value;
    $("start").disabled = value;
}

function setProgress(value, total) {
    let progress = $("progress");
    progress.min = 0;
    progress.max = total;
    progress.value = value;
}