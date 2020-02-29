
importScripts('https://www.gstatic.com/firebasejs/6.3.4/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/6.3.4/firebase-messaging.js');

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

// Initialize the Firebase app in the service worker by passing in the
// messagingSenderId.
firebase.initializeApp(firebaseConfig);
  
  // Retrieve an instance of Firebase Messaging so that it can handle background
  // messages.
const messaging = firebase.messaging();

