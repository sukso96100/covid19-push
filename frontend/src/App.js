import React, {useEffect} from 'react';
import logo from './logo.svg';
import './App.css';
import * as firebase from "firebase/app";
import "firebase/messaging";
import localForage from "localforage";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Redirect from './redirect';

const firebaseConfig = {
  apiKey: "AIzaSyCJC0XNjwo_HUKpH1FwSxYQAxlF3O-Uzes",
  authDomain: "covid19-269505.firebaseapp.com",
  databaseURL: "https://covid19-269505.firebaseio.com",
  projectId: "covid19-269505",
  storageBucket: "covid19-269505.appspot.com",
  messagingSenderId: "649845923341",
  appId: "1:649845923341:web:5dd4d71ca9ec0daa383f44",
  measurementId: "G-KYQSGTFKNG"
};
const vapid = "BHLL-JhczA92RQd2uLIaAFEqICgVFapkcxaRsxUC_p2E_bfaftWKolfM7rgx2jxCH3IPbe3jRIbudxzf0frC5N0"


firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();
messaging.usePublicVapidKey(vapid); 

export default function App() {
  return (
    <Router>
      <div>

        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/redirect/:url">
            <Redirect />
          </Route>
          <Route path="/">
            <Home />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

function Home() {
  useEffect(()=>{
    messaging.onMessage((payload) => {
    console.log('Message received. ', payload);
            // ...
    });
    messaging.onTokenRefresh(() => {
      messaging.getToken().then((refreshedToken) => {
        localForage.setItem("token", refreshedToken)
        localForage.setItem("tokenSent", "1")
      }).catch((err) => {
        console.log('Unable to retrieve refreshed token ', err);
      });
    });
  },[])
  const subscribe = async() => {
    let result = await Notification.requestPermission();
        if(result=="granted"){
          let token = await messaging.getToken();
          

          if(!tokenSaved()){
            localForage.setItem("token", token)
            localForage.setItem("tokenSent", "1")
          }
          subscribePush(token);
          
        }else{
          alert("알림 권한을 승인해야 알림을 수신할 수 있습니다.")
        }
  }
  return (
    <div class="main">
      <h1 class="title">코로나19 알리미</h1>
      <button onClick={subscribe}>알림 받기</button>
      <div class="box stat">
          <div class="statitem">
              <span>0</span>
              <b>확진</b>
          </div>
          <div class="statitem">
              <span>0</span>
              <b>완치</b>
          </div>
          <div class="statitem">
              <span>0</span>
              <b>사망</b>
          </div>
      </div>
      <strong>Messages</strong>
      <br/>
      <div id="messages"></div>
  </div>
  );
}


async function tokenSaved(){
  let token = await localForage.getItem("token");
  return token != "";
}

function subscribePush(token){
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

function unsubscribePush(token){
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