export async function getStat(){
    let result = await fetch("/stat");
    return result.json()
}

export async function getNews(){
    let result = await fetch("/news");
    return result.json()
}

export function subscribePush(token){
    fetch("/subscribe/stat",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    fetch("/subscribe/news",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
  }
  
  export function unsubscribePush(token){
    fetch("/unsubscribe/stat",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
    fetch("/unsubscribe/news",{
      method: 'POST', // or 'PUT'
      body: JSON.stringify({
        "token": token
      }), // data can be `string` or {object}!
      headers:{
        'Content-Type': 'application/json'
      }
    })
  }