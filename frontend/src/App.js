import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
  import * as firebase from "firebase/app";
import "firebase/messaging";
import localForage from "localforage";
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import Redirect from './redirect';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Divider from '@material-ui/core/Divider';
import ListItemText from '@material-ui/core/ListItemText';
import {subscribePush, unsubscribePush, getNews, getStat} from './api';
import {firebaseConfig, vapidKey} from './fcmconfig';


const useStyles = makeStyles({
  root:{
    padding: '16px'
  },
  subBtns: {
    padding: '8px',
    margin: '8px'
  },
  stat: {
    display: 'flex',
    flexDirection: 'row'
  },
  statitem: {
    flex:1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  ststdata: {
    fontSize: '64'
  },
  card:{
    margin: '8px'
  }
});



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

firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();
messaging.usePublicVapidKey(vapidKey); 

function Home() {
  const classes = useStyles();
  let [isSubscribed, setSubscribed] = useState(false)
  let [statData, setStatData] = useState({confirmed:0, cured:0, death:0})
  let [newsData, setNewsData] = useState([])
  useEffect(()=>{
    (async function(){
      setSubscribed(await tokenSaved())
      setStatData(await getStat())
      setNewsData(await getNews())
    })()
    messaging.onMessage(async(payload) => {
         const {title, ...options} = payload.notification;
        navigator.serviceWorker.ready.then(registration => {
            registration.showNotification(title, options);
        });
        setStatData(await getStat())
        setNewsData(await getNews())
    });
    messaging.onTokenRefresh(() => {
      messaging.getToken().then(async(refreshedToken) => {
        await localForage.setItem("token", refreshedToken)
        await localForage.setItem("tokenSent", "1")
      }).catch((err) => {
        console.log('Unable to retrieve refreshed token ', err);
      });
    });
  },[])
  const subscribe = async() => {
    let result = await Notification.requestPermission();
        if(result=="granted"){
          let token = await messaging.getToken();
          await localForage.setItem("token", token)
          await localForage.setItem("tokenSent", "1")
          subscribePush(token);
          setSubscribed(await tokenSaved())
        }else{
          alert("알림 권한을 승인해야 알림을 수신할 수 있습니다.")
        }
  }
  const unsubscribe = async()=>{
    await localForage.setItem("token", "")
    await localForage.setItem("tokenSent", "0")
    let token = await messaging.getToken();
    unsubscribePush(token);
    console.log(await tokenSaved())
    setSubscribed(await tokenSaved())
  }
  return (
    <div className={classes.root}>
      <h1 class="title">코로나19 알리미</h1>
      <p>질병관리본부 코로나19 홈페이지에서 발생 동향과 새 공지사항을 푸시알림으로 알려드립니다.</p>
      <b>Web Notification 기능을 지원하는 웹 브라우저에서 알림 권한 허용 후 이용 가능합니다.</b><br/>
      <Button variant="contained" color="primary" className={classes.subBtns} onClick={subscribe}>
        알림 구독
      </Button>
      <Button variant="contained" color="primary" className={classes.subBtns} onClick={unsubscribe}>
        구독 해제
      </Button>
      <Card className={classes.card}>
      <CardContent>
        <Typography color="textSecondary" gutterBottom>
          코로나19 발생 현황
        </Typography>
        <div className={classes.stat}>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.confirmed}</Typography>
              <b>확진</b>
          </div>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.cured}</Typography>
              <b>완치</b>
          </div>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.death}</Typography>
              <b>사망</b>
          </div>
      </div>
      </CardContent>
      <CardActions>
        <Button size="small" href="http://ncov.mohw.go.kr/bdBoardList_Real.do">자세히 보기</Button>
      </CardActions>
    </Card>
    <Card className={classes.card}>
      <CardContent>
    <Typography color="textSecondary" gutterBottom>
          질병관리본부 공지사항
          </Typography>
    <List>
      {newsData.map((item, i)=>(
        <div>
          <ListItem alignItems="flex-start" button
            onClick={()=>{
              window.open(item.link, "_blank")
            }}>
          <ListItemText
            primary={item.title}
            secondary={
              <React.Fragment>
                <Typography
                  component="span"
                  variant="body2"
                  className={classes.inline}
                  color="textPrimary"
                >
                  {item.dept}
                </Typography>
                {/* {" — I'll be in your neighborhood doing errands this…"} */}
              </React.Fragment>
            }
          />
        </ListItem>
        <Divider  component="li" />
        </div>
      ))}
      </List>
      </CardContent>
      </Card>
      <a href="https://youngbin.xyz">개발자 개인 웹사이트 방문</a><br/>
      <a href="mailto:sukso96100@gmail.com">개발자와 연락하기(이메일)</a>
  </div>
  );
}


async function tokenSaved(){
  let token = await localForage.getItem("token");
  return token != undefined && token === "";
}
