const firebaseConfig = {
    apiKey: process.env.REACT_APP_FIREBASE_APIKEY, 
    authDomain: process.env.REACT_APP_FIREBASE_AUTHDOMAIN,
    databaseURL: process.env.REACT_APP_FIREBASE_DBURL,
    projectId: process.env.REACT_APP_FIREBASE_PROJID,
    storageBucket: process.env.REACT_APP_FIREBASE_BUCKET,
    messagingSenderId: process.env.REACT_APP_FIREBASE_SENDERID,
    appId: process.env.REACT_APP_FIREBASE_APPID,
    measurementId:process.env.REACT_APP_FIREBASE_ANALYTICS
  };
const vapidKey = process.env.REACT_APP_FIREBASE_VAPIDKEY
export {firebaseConfig, vapidKey}