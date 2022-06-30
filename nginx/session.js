const maxQueryCountPerExpire = 2

function handleRequest(r) {
    let sessionToken = ""
    let captchaKey = ""
    const cookies = r.headersIn['Cookie'].replace(/ /g,'').split(';')
    for (let i = 0; i < cookies.length; i++){
        const cookie = cookies[i].split('=')

        if (cookie[0] === 'captcha-key'){
            captchaKey = cookie[1]
        }
        if (cookie[0] === 'session-token'){
            sessionToken = cookie[1]
        }
    }

    adapterQuery(r,'get ' + sessionToken + '-sessionKey').then(function (response) {
            const resp = response.responseBody.split('\n')[1].trim();
            if (resp === '0' || resp === 0){
                let params = JSON.stringify({
                    "query_count": 1,
                    "expire_at": new Date(new Date().getTime() + 10000),
                });

                return adapterQuery(r,"set " + sessionToken + "-sessionKey '" + params + "'").catch(function (reason) {
                    throw new Error("error setting new session token " + reason);
                })
            }

            const sessionParamsMap = JSON.parse(resp);
            const expireDate = new Date(sessionParamsMap['expire_at']);
            const queryCount = sessionParamsMap['query_count'];

            if (Date.now() >= expireDate){
                let params = JSON.stringify({
                    "query_count": 1,
                    "expire_at": new Date(new Date().getTime() + 10000),
                });

                return adapterQuery(r,"set " + sessionToken + "-sessionKey '" + params + "'").catch(function(reason) {
                    throw new Error("error setting new session token " + reason);
                });
            }else{
                if (queryCount < maxQueryCountPerExpire){
                    let params = JSON.stringify({
                        "query_count": queryCount + 1,
                        "expire_at": new Date(new Date().getTime() + 10000),
                    });

                    return adapterQuery(r,"set " + sessionToken + "-sessionKey '" + params + "'").catch(function (reason) {
                        throw new Error("error setting new session token " + reason);
                    })
                }else{
                    throw new Error('session limit reached');
                }
            }
    }).then(ans =>{
        r.internalRedirect('@home');
    }).catch( reason => {
        if (captchaKey != null && captchaKey !== '' && reason.message === 'session limit reached'){
            r.headersOut['captcha-key'] = captchaKey;
            r.internalRedirect('@captcha');
        }else{
            r.subrequest('/generate_session', {
                args: '',
                body: '',
                method: 'GET'
            }, function (sessionResponse) {
                const parsedSessionResponse = JSON.parse(sessionResponse.responseBody);
                const captchaKey = parsedSessionResponse['captcha_key']
                const expTime = parsedSessionResponse['exp_time']

                r.headersOut['Set-Cookie'] = [
                    'captcha-key='+captchaKey,
                    'exp-time='+expTime
                ];
                r.headersOut['Content-Type'] = 'text/html';
                r.headersOut['captcha-key'] = captchaKey;
                r.headersOut['exp-time'] = expTime;
                r.internalRedirect('@captcha');
            })
        }
    });
}

function adapterQuery(r, query) {
    return r.subrequest('/redisadapter', `query=${query}\r\n`)
}

export default {
    handleRequest
};