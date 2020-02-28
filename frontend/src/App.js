import React, {useEffect} from 'react';
import logo from './logo.svg';
import './App.css';
import * as firebase from "firebase/app";
import "firebase/messaging";

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


firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();

function App() {
  useEffect(()=>{
    (async function(){
      let result
        result = await Notification.requestPermission();
        if(result=="granted"){
          token = await messaging.getToken();
        }else{
          alert("알림 권한을 승인해야 알림을 수신할 수 있습니다.")
        }
    });
  },[]);
  return (
    <div class="main">
      <h1 class="title">코로나19 알리미</h1>
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

export default App;
