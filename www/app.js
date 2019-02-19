(function (window, undefined) {
    'use strict'

    class CanvasPainter extends HTMLElement {
        constructor() {
            super();
            const styles = `
            <style>
                .info {
                    font-weight: 300;
                    margin-top: 10px;
                }
            </style>
            `;

            const template = `
                <h1 class="title">${this.getAttribute('text')}<h1>
                <canvas class="canvas-area" id="canvas" width="280" height="280"></canvas>
                <div class="info" id="info"></div>
            `;

            this._shadow = this.attachShadow({ mode: 'open' });
            this._shadow.innerHTML = styles + template;

            const canvas = this._canvas;
            canvas.setAttribute('width', 280);
            canvas.setAttribute('height', 280);
            canvas.addEventListener('mousedown', this._mousedown);
            canvas.addEventListener('mousemove', this._mousemove);
            canvas.addEventListener('mouseup', this._mouseup);
            canvas.addEventListener('mouseup', this._mouseout);

            const ctx = canvas.getContext('2d');
            ctx.lineJoin = 'round';
            ctx.lineCap = 'round';
            ctx.lineWidth = 25;
            ctx.strokeStyle = '#fff';
            ctx.fillStyle = '#000';
            ctx.fillRect(0, 0, canvas.width, canvas.height);


            this._ctx = ctx;
            this._lastY = 0;
            this._lastX = 0;
            this._isDrawing = false;

        }

        getByteArray() {
            const input = this._canvas;
            const output = document.createElement('canvas');
            output.width = 28;
            output.height = 28;

            const ctx = output.getContext("2d");
            ctx.scale(.1, .1)
            ctx.drawImage(input, 0, 0);

            return this._b64ToUint8Array(output.toDataURL('image/png'));
        }

        setInfo(text) {
            this._info.textContent = text;
        }

        get _canvas() {
            return this._shadow.getElementById('canvas');
        }

        get _info() {
            return this._shadow.getElementById('info');
        }

        _mousedown = (e) => {
            if (e.which !== 1) { return; }

            this._isDrawing = true;
            this._lastX = e.offsetX;
            this._lastY = e.offsetY;
        }

        _mousemove = (e) => {
            if (!this._isDrawing) { return; }

            this._ctx.beginPath();
            this._ctx.moveTo(this._lastX, this._lastY);
            this._ctx.lineTo(e.offsetX, e.offsetY);
            this._ctx.stroke();

            this._lastX = e.offsetX;
            this._lastY = e.offsetY;
        }

        _mouseup = () => {
            this._isDrawing = false;
        }

        _mouseout = () => {
            this._isDrawing = false;
        }

        _b64ToUint8Array(b64Image) {
            const img = window.atob(b64Image.split(',')[1]);
            const buffer = [];
            for (let i = 0, len = img.length; i < len; i++) {
                buffer.push(img.charCodeAt(i));
            }
            return new Uint8Array(buffer);
        }
    }

    class NeuralCalculator extends HTMLElement {
        constructor() {
            super();

            const styles = `
                <style>
                    .calculator {               
                        display: flex;
                        align-items: center;
                        justify-content: space-around;                
                        padding: 2em;
                        padding-bottom: 3em;
                        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.2), 0 25px 50px 0 rgba(0, 0, 0, 0.1);
                    }
        
                    .input,
                    .output {
                        border: 1px solid #ccc
                    }
        
                    .output {
                        width: 280px;
                        height: 280px;
                        text-align: center;
                        line-height: 280px;
                        font-size: 200px;
                    }
        
                    .send {
                        line-height: 40px;
                        font-size: 2em;
                        margin-top: 10px;
                        box-shadow: 0 3px 1px -2px rgba(0, 0, 0, .2), 0 2px 2px 0 rgba(0, 0, 0, .14), 0 1px 5px 0 rgba(0, 0, 0, .12);
                        box-sizing: border-box;
                        position: relative;
                        -webkit-user-select: none;
                        -moz-user-select: none;
                        -ms-user-select: none;
                        user-select: none;
                        cursor: pointer;
                        outline: 0;
                        border: none;
                        -webkit-tap-highlight-color: transparent;
                        display: inline-block;
                        white-space: nowrap;
                        text-decoration: none;
                        vertical-align: baseline;
                        text-align: center;
                        min-width: 64px;
                        line-height: 36px;
                        padding: 5px 16px;
                        border-radius: 4px;
                        overflow: visible;
                        transform: translate3d(0, 0, 0);
                        transition: background .4s cubic-bezier(.25, .8, .25, 1), box-shadow 280ms cubic-bezier(.4, 0, .2, 1);
                    }
        
                    .send:active {
                        box-shadow: 0 5px 5px -3px rgba(0, 0, 0, .2), 0 8px 10px 1px rgba(0, 0, 0, .14), 0 3px 14px 2px rgba(0, 0, 0, .12);
                    }
        
                    .center {
                        text-align: center
                    }
                </style>
            `

            const template = `
                <div class="calculator">
                    <canvas-painter id="input1" text="Draw 0-9" class="center"></canvas-painter>
                    <h1 class="operator">X</h1>
                    <canvas-painter id="input2" text="Draw 0-9" class="center"></canvas-painter>
                    <button type="button" id="send" class="send">=</button>
                    <div class="center">
                        <h1 class="res-title">Result</h1>
                        <div id="output" class="output"></div>
                    </div>
                </div>
            `

            this._shadow = this.attachShadow({ mode: 'open' });
            this._shadow.innerHTML = styles + template;

            this._btn.addEventListener('click', this.calc);
        }

        calc = async () => {
            const results = await Promise.all([
                this._sendImage(this._input1.getByteArray()),
                this._sendImage(this._input2.getByteArray())
            ]);

            this._showResult(results);
        }

        get _input1() {
            return this._shadow.getElementById('input1');
        }

        get _input2() {
            return this._shadow.getElementById('input2');
        }

        get _output() {
            return this._shadow.getElementById('output');
        }

        get _btn() {
            return this._shadow.getElementById('send');
        }

        _showResult(results) {
            const a1 = results[0].answer;
            const a2 = results[1].answer;
            
            const result = a1 * a2;

            this._input1.setInfo(a1);
            this._input2.setInfo(a2);
            this._output.textContent = result;
        }

        _sendImage(byteArray) {
            const formData = new FormData();
            formData.append('file', new Blob([byteArray], { type: 'image/png' }));

            return new Promise((resolve, reject) => {
                const xhr = new XMLHttpRequest();
                xhr.open('POST', this.getAttribute('url'));

                xhr.onreadystatechange = () => {
                    if (xhr.readyState === 4) {
                        if (xhr.status === 200) {
                            resolve(JSON.parse(xhr.response));
                        } else {
                            reject(JSON.parse(xhr.response));
                        }
                    }
                };

                xhr.send(formData);
            });
        }

    }

    customElements.define('canvas-painter', CanvasPainter);
    customElements.define('neural-calculator', NeuralCalculator);

})(window);