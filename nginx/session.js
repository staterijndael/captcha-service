const maxQueryCountPerExpire = 2

function handleRequest(r) {
    adapterQuery(r,'get ' + r.headersIn['session-token'] + '-sessionKey').then(function (response) {
            const resp = response.responseBody.split('\n')[1].trim();
            if (resp === '0' || resp === 0){
                let params = JSON.stringify({
                    "query_count": 1,
                    "expire_at": new Date(new Date().getTime() + 10000),
                });

                return adapterQuery(r,"set " + r.headersIn['session-token'] + "-sessionKey '" + params + "'").catch(function (reason) {
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

                return adapterQuery(r,"set " + r.headersIn['session-token'] + "-sessionKey '" + params + "'").catch(function(reason) {
                    throw new Error("error setting new session token " + reason);
                });
            }else{
                if (queryCount < maxQueryCountPerExpire){
                    let params = JSON.stringify({
                        "query_count": queryCount + 1,
                        "expire_at": new Date(new Date().getTime() + 10000),
                    });

                    return adapterQuery(r,"set " + r.headersIn['session-token'] + "-sessionKey '" + params + "'").catch(function (reason) {
                        throw new Error("error setting new session token " + reason);
                    })
                }else{
                    throw new Error('session limit reached');
                }
            }
    }).then(ans =>{
        r.subrequest('/home').then(function(response){
            r.return(200, response.responseBody)
        });
    }).catch(function (reason){
        r.return(400, reason);
    });
}

function adapterQuery(r, query) {
    return r.subrequest('/redisadapter', `query=${query}\r\n`)
}

export default {
    handleRequest
};