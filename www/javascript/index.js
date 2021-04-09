if (! WebAssembly.instantiateStreaming){
	 
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();

let mod, inst;

WebAssembly.instantiateStreaming(fetch("/wasm/update_exif.wasm"), go.importObject).then(
    
    async result => {
	document.getElementById("button").innerText = "Update";
	document.getElementById("button").removeAttribute("disabled");
        mod = result.module;
        inst = result.instance;
	await go.run(inst);
    }
);

async function update() {

    var img = document.getElementById("image");
    
    var canvas = document.createElement("canvas");
    canvas.width = img.width;
    canvas.height = img.height;
    var ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    var b64_img = canvas.toDataURL("image/jpeg", 1.0);

    var update = { "CameraOwnerName": "Bob" };
    var enc_update = JSON.stringify(update);
    
    var rsp = update_exif(b64_img, enc_update);

    if (! rsp){
	return;
    }
    
    var blob = dataURLToBlob(rsp);

    if (! blob){
	return;
    }
    
    saveAs(blob, "example.jpg");
}

var dataURLToBlob = function(dataURL){
    var BASE64_MARKER = ";base64,";
    if (dataURL.indexOf(BASE64_MARKER) == -1)
    {
        var parts = dataURL.split(",");
        var contentType = parts[0].split(":")[1];
        var raw = decodeURIComponent(parts[1]);

        return new Blob([raw], {type: contentType});
    }

    var parts = dataURL.split(BASE64_MARKER);
    var contentType = parts[0].split(":")[1];
    var raw = window.atob(parts[1]);
    var rawLength = raw.length;

    var uInt8Array = new Uint8Array(rawLength);

    for (var i = 0; i < rawLength; ++i) {
        uInt8Array[i] = raw.charCodeAt(i);
    }

    return new Blob([uInt8Array], {type: contentType});
}
