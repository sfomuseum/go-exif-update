var sfomuseum = sfomuseum || {};

sfomuseum.wasm = (function(){

    var self = {

	fetch: function(wasm_uri){

	    var pending = 1;
	    
	    return new Promise((resolve, reject) => {
		
		if (! WebAssembly.instantiateStreaming){
		    
		    WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		    };
		}
		
		const export_go = new Go();
		
		let export_mod, export_inst;	

		// See this, with the headers? This is important if we're running in
		// a AWS Lambda + API Gateway context. Without this API Gateway will
		// return the WASM binary as a base64-encoded blob. Note that this
		// also depends on configuring both the API Gateway and the 'lambda://'
		// server URI to specify that 'application/wasm' is treated as binary
		// data. Computers, amirite...
		    
		var fetch_headers = new Headers();
		fetch_headers.set("Accept", "application/wasm");
		
		const fetch_opts = {
		    headers: fetch_headers,
		};

		self.log("fetch " + wasm_uri);
		
		WebAssembly.instantiateStreaming(fetch(wasm_uri, fetch_opts), export_go.importObject).then(
		    
		    async result => {

			self.log("retrieved " + wasm_uri);			
			
			pending -= 1;
			
			if (pending == 0){
			    resolve();
			}

			export_mod = result.module;
			export_inst = result.instance;
			await export_go.run(export_inst);
		    }
		);
		
	    });
	},

	'log': function(msg) {
	    var dt = new Date();
	    console.log("[wasm][" + dt.toISOString() + "] fetch " + msg);	    
	}
    };

    return self;
})();
