import React, {useEffect, useState} from 'react';

export function CompatInfo(){
    const [hasSW, setHasSW] = useState(true);
    const [hasNoti, setHasNoti] = useState(true);
    const [hasPush, setHasPush] = useState(true);
    useEffect(()=>{
      setHasSW('serviceWorker' in navigator)
      setHasNoti("Notification" in window)
      setHasPush('PushManager' in window)
    },[])
    return(
      <div>
        {(!hasSW || !hasNoti)?
        (<b>
          웹 브라우저가 다음 기능을 제공하지 않아 알림 구독이 불가능합니다.<br/>
          Telegram 채널 구독을 대신 이용해 주세요<br/></b>):(
        <b>알림 권한 허용 후 이용해 주세요.</b>)}
        {hasSW?(<b></b>):(<b>→ 서비스워커(Service Worker)<br/></b>)}
        {hasNoti?(<b></b>):(<b>→ 웹 알림(Web Notification)<br/></b>)}
        {hasPush?(<b></b>):(<b>→ 웹 푸시(Web Push)<br/></b>)}
      </div>
    )
  }

export function isSupported(){
    return ('serviceWorker' in navigator) 
        && ("Notification" in window)
        && ("PushManager" in window)
}