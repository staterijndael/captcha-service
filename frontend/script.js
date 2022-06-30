const sessionBackendAddr = "http://0.0.0.0:8082"

async function renderBoxImage() {
    try {
        let storagedCaptchaKey = localStorage.getItem('captcha_key')
        let storagedSessionToken = localStorage.getItem('session_token')
        if (storagedCaptchaKey == null && storagedSessionToken == null) {
            storagedCaptchaKey = await generateSession();
        }

        if (storagedCaptchaKey != null && storagedSessionToken == null) {
            let captcha_text = document.getElementById('captcha_text').value
            if (captcha_text === null || captcha_text === '') {
                hideErrorMsg()
                renderCaptchaImage(storagedCaptchaKey)
            }else{
                try {
                    let decryptedText = decryptAES128(captcha_text, storagedCaptchaKey)
                    if (decryptedText.length < captcha_text.length || captcha_text !== decryptedText.substring(0, captcha_text.length)) {
                        hideCaptchaImage()
                        printErrMessage("Вы шото не то написали")
                    }

                    let sessionToken = ""
                    for (let i = captcha_text.length; decryptedText[i] !== '0' && i < decryptedText.length; i++){
                        sessionToken += decryptedText[i]
                    }

                    localStorage.setItem('session_token', sessionToken)
                    storagedSessionToken = sessionToken
                }catch (e){
                    hideCaptchaImage()
                    printErrMessage("Вы шото не то написали")
                }
            }
        }

        if (storagedCaptchaKey != null && storagedSessionToken != null){
            try {
                hideErrorMsg()
                hideCaptchaImage()
                await renderHomeImage(storagedSessionToken)
            }catch (e){
                console.log(e.responseBody)
            }
        }
    } catch (e) {
        console.log(e);
    }
}

function decryptAES128(key, cipherText){
    let sha256Key = CryptoJS.SHA256(key).toString();
    sha256Key = sha256Key.substring(0, sha256Key.length/2);

    // Fix 3: disable padding
    const decrypted = CryptoJS.AES.decrypt(
        {ciphertext: CryptoJS.enc.Hex.parse(cipherText)}, // Fix 1: pass a CipherParams object
        CryptoJS.enc.Utf8.parse(sha256Key), // Fix2: UTF-8 encode the key
        {mode: CryptoJS.mode.ECB, padding: CryptoJS.pad.NoPadding})
    return decrypted.toString(CryptoJS.enc.Utf8)
}

function renderCaptchaImage(captchaKey){
    document.getElementById('boxImage').src = sessionBackendAddr + "/img?captcha_key=" + captchaKey

    let captchaTextElements = document.getElementsByClassName('captcha_text');
    for (let i = 0; i < captchaTextElements.length; i++){
        captchaTextElements[i].style.visibility = "visible";
    }
}

async function renderHomeImage(sessionToken){
    const src = sessionBackendAddr + "/";
    const options = {
        method: 'GET',
        headers: {
            'session-token': sessionToken,
            'Cache-Control': 'no-cache',
        },
    };

    console.log(sessionToken)
    fetch(src, options)
        .then(res => {return res.text()})
        .then(html => {
            const homeImage = document.getElementById('homeImage')
            homeImage.innerHTML = html;
            homeImage.style.visibility = 'visible';
        });
}

function hideCaptchaImage(){
    document.getElementById('boxImage').src = ""

    document.getElementById('captcha_text').value = ""

    let captchaTextElements = document.getElementsByClassName('captcha_text');
    for (let i = 0; i < captchaTextElements.length; i++){
        captchaTextElements[i].style.visibility = "hidden";
    }
}

function printErrMessage(msg){
    let errorMsgElem = document.getElementById('error_msg');
    errorMsgElem.style.visibility = "visible";
    errorMsgElem.innerHTML = msg;
}

function hideErrorMsg(){
    let errorMsgElem = document.getElementById('error_msg');
    errorMsgElem.style.visibility = "hidden";
    errorMsgElem.innerHTML = "";
}

async function generateSession(){
    const Http = new XMLHttpRequest();
    const url=sessionBackendAddr + "/generate_session";
    Http.open("GET", url);
    Http.send();

    const response = await fetch(sessionBackendAddr + "/generate_session", {}) // type: Promise<Response>
    if (!response.ok) {
        throw 'got error during generating session ' + response.statusText;
    }

    const resp = JSON.parse(Http.responseText)
    localStorage.setItem('captcha_key', resp['captcha_key'])

    return resp['captcha_key'];
}