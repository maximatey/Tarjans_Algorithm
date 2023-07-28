const goWasm = new Go();
let instance;

WebAssembly.instantiateStreaming(fetch("maingo.wasm"), goWasm.importObject)
    .then((result) => {
        instance = result.instance;
        goWasm.run(instance);

        document.getElementById("run-function").addEventListener("click", () => {
            const inputText = document.getElementById("function-input").value;
            console.log("1")
            const outputDiv = document.getElementById("output");
            console.log("2")
            const outstr = mainFunc(inputText);
            console.log(outstr)

            outputDiv.innerHTML = `<pre>${outstr}</pre>`
        

            // if (typeof mainFunc === "function") {
            //     const result = mainFunc(inputText); // Call the mainFunc with the textarea content
            //     outputDiv.innerHTML = `<pre>${result}</pre>`; // Display the output in a <pre> tag for better formatting
            // } else {
            //     outputDiv.innerHTML = `<p>Function 'mainFunc' not found.</p>`;
            // }
        });

        // document.getElementById("run-function2").addEventListener("click", () => {
        //     const inputText = document.getElementById("file-selector").files[0];
        //     console.log("1")
        //     const outputDiv = document.getElementById("output");
        //     console.log("2")

        //     if (inputText) {
        //         const outstr = mainFunc2(inputText);
        //         console.log(outstr)
    
        //         outputDiv.innerHTML = `<pre>${outstr}</pre>`
        //     } else {
        //         outputDiv.innerHTML = `<pre>No File Detected</pre>`
        //     }
            
        

        //     // if (typeof mainFunc === "function") {
        //     //     const result = mainFunc(inputText); // Call the mainFunc with the textarea content
        //     //     outputDiv.innerHTML = `<pre>${result}</pre>`; // Display the output in a <pre> tag for better formatting
        //     // } else {
        //     //     outputDiv.innerHTML = `<p>Function 'mainFunc' not found.</p>`;
        //     // }
        // });
    })
    .catch((err) => {
        console.error("Error while loading WebAssembly:", err);
    });
