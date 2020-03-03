export async function getStat(){
    let result = await fetch("/stat");
    return result.json()
}

export async function getNews(){
    let result = await fetch("/news");
    return result.json()
}

export async function subscribePush(token){
    let r1 = await fetch("/subscribe/stat",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    let r2 = await fetch("/subscribe/news",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    return r1.ok && r2.ok
  }
  
  export async function unsubscribePush(token){
    let r1 = await fetch("/unsubscribe/stat",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    let r2 = await fetch("/unsubscribe/news",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    return r1.ok && r2.ok
  }