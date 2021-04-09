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

	var update_button = document.getElementById("update");
	var add_button = document.getElementById("add");    

	update_button.innerText = "Update";
	update_button.removeAttribute("disabled");
	update_button.onclick = update;

	add_button.innerText = "Add Property";	
	add_button.removeAttribute("disabled");
	add_button.onclick = add_property;
	
        mod = result.module;
        inst = result.instance;
	await go.run(inst);
    }
);

async function add_property(){

    var props = document.getElementsByClassName("exif-property");
    var count = props.length;

    var idx = count + 1;
    var id = "exif-property-" + idx;

    var t_id = "exif-property-tag" + idx;
    var v_id = "exif-property-value" + idx;        

    var group = document.createElement("div");
    group.setAttribute("class", "form-group exif-property");
    group.setAttribute("data-index", idx);    
    group.setAttribute("id", id);

    var input_t = document.createElement("input");
    input_t.setAttribute("type", "input");
    input_t.setAttribute("placeholder", "A valid EXIF tag name");
    input_t.setAttribute("id", t_id);

    var input_v = document.createElement("input");
    input_v.setAttribute("type", "input");
    input_v.setAttribute("placeholder", "A valid EXIF tag value");
    input_v.setAttribute("id", v_id);

    group.appendChild(input_t);
    group.appendChild(input_v);    
    
    var form = document.getElementById("properties-form");    
    form.appendChild(group);
    
    var update_button = document.getElementById("update");
    update_button.style.display = "block";
    
    return false;
}

async function update() {

    var props = document.getElementsByClassName("exif-property");
    var count = props.length;

    if (count == 0){
	return false;
    }

    var update = {};
    var has_updates = false;
    
    for (var i=0; i < count; i++){

	var el = props[i];
	var idx = el.getAttribute("data-index");

	var t_id = "exif-property-tag" + idx;
	var v_id = "exif-property-value" + idx;

	var t_el = document.getElementById(t_id);
	var v_el = document.getElementById(v_id);	

	var t = t_el.value;
	var v = v_el.value;

	if (t == ""){
	    continue;
	}

	if (v == ""){
	    continue;
	}

	update[t] = v;
	has_updates = true;
    }

    if (! has_updates){
	return false;
    }

    console.log(update);    
    var enc_update;
    
    try {
	enc_update = JSON.stringify(update);
    } catch(e){
	console.log("Failed to marhsal update", update, e);
	return false;
    }
    
    var img = document.getElementById("image");
    
    var canvas = document.createElement("canvas");
    canvas.width = img.width;
    canvas.height = img.height;
    var ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    var b64_img = canvas.toDataURL("image/jpeg", 1.0);
    
    var rsp = update_exif(b64_img, enc_update);

    if (! rsp){
	return false;
    }
    
    var blob = dataURLToBlob(rsp);

    if (! blob){
	return false;
    }
    
    saveAs(blob, "example.jpg");
    return false;
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
