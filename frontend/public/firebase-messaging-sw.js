
importScripts('https://www.gstatic.com/firebasejs/6.3.4/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/6.3.4/firebase-messaging.js');

const firebaseConfig = {
    messagingSenderId: "649845923341"
  };

// Initialize the Firebase app in the service worker by passing in the
// messagingSenderId.
firebase.initializeApp(firebaseConfig);
  
  // Retrieve an instance of Firebase Messaging so that it can handle background
  // messages.
const messaging = firebase.messaging();
messaging.setBackgroundMessageHandler((payload)=>{
  console.log(payload)
})
