import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import * as firebase from "firebase/app";
import 'firebase/analytics';
import "firebase/messaging";
import localForage from "localforage";
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import {CompatInfo, isSupported} from './compatCheck';
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
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import LanguageIcon from '@material-ui/icons/Language';
import AlternateEmailIcon from '@material-ui/icons/AlternateEmail';


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
    flexDirection: 'row',
    marginTop: '16px',
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
  },
  iconBtns:{
    display: 'flex',
    flexDirection: 'row'
  },
  iconBtnsItem:{
    padding: '8px',
    margin: '4px',
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
            <Redirect ga={analytics}/>
          </Route>
          <Route path="/">
            <Home/>
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

firebase.initializeApp(firebaseConfig);
const messaging = isSupported() ? firebase.messaging() : null;
const analytics = firebase.analytics();

if(isSupported()){
  messaging.usePublicVapidKey(vapidKey); 
}

function Home(props) {
  const classes = useStyles();
  const [snackbar, setSnackbar] = React.useState(false);
  const [snackMsg, setSnackMsg] = React.useState('');
  const [statData, setStatData] = useState({
    confirmed:0, cured:0, death:0, checking:0, patients:0, resultNeg:0
  })
  let [newsData, setNewsData] = useState([])
  useEffect(()=>{
    analytics.logEvent('open_homepage');

    (async function(){
      setStatData(await getStat())
      setNewsData(await getNews())
    })()
    if(isSupported()){
      messaging.onMessage(async(payload) => {
        console.log('msg received')
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
    }
  },[])

  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbar(false);
  };
  const subscribe = async() => {
    if(isSupported()){
      let result = await Notification.requestPermission();
      if(result=="granted"){
        let token = await messaging.getToken();
        await localForage.setItem("token", token)
        await localForage.setItem("tokenSent", "1")
        let result = await subscribePush(token);
        if(result){
          setSnackMsg('알림 구독 완료.');
          analytics.logEvent('subscribe', {result: 'ok'});
        }else{
          setSnackMsg('알림 구독중 오류가 발생했습니다.');
          analytics.logEvent('subscribe', {result: 'error'});
        }
        
        setSnackbar(true);
      }else{
        setSnackMsg('알림 권한을 허용해야 이용하실 수 있습니다.');
        setSnackbar(true);
        analytics.logEvent('subscribe', {result: 'no_permission'});
      }
    }else{
      setSnackMsg('사용중인 웹 브라우저에서 이용하실 수 없습니다.');
      setSnackbar(true);
      analytics.logEvent('subscribe', {result: 'unsuported'});
    }
  }
  const unsubscribe = async()=>{
    if(isSupported()){
      await localForage.setItem("token", "")
      await localForage.setItem("tokenSent", "0")
      let token = await messaging.getToken();
      let result = await unsubscribePush(token);
      if(result){
        setSnackMsg('알림 구독 해제 완료.');
        analytics.logEvent('unsubscribe', {result: 'ok'});
      }else{
        setSnackMsg('알림 구독 해제중 오류가 발생했습니다..');
        analytics.logEvent('unsubscribe', {result: 'error'});
      }
      setSnackbar(true);
    }else{
      setSnackMsg('사용중인 웹 브라우저에서 이용하실 수 없습니다.');
      setSnackbar(true);
      analytics.logEvent('subscribe', {result: 'unsupported'});
    }
  }
  return (
    <div className={classes.root}>
      <h1 class="title">코로나19 푸시알리미</h1>
      
      <p>질병관리본부 코로나19 홈페이지에서 발생 동향과 새 공지사항을 푸시알림으로 알려드립니다.</p>
      <CompatInfo/>
     <br/>
      <Button variant="outlined" color="primary" className={classes.subBtns} onClick={subscribe}>
        알림 구독
      </Button>
      <Button variant="outlined" color="primary" className={classes.subBtns} onClick={unsubscribe}>
        구독 해제
      </Button>
      <Button variant="outlined" color="primary" className={classes.subBtns}
        onCliek={()=>{
          window.open("https://t.me/covid19push", "_blank")
          analytics.logEvent('telegram');
        }}>
        Telegram 채널 구독
      </Button>
      <Card className={classes.card}>
      <CardContent>
        <Typography color="textSecondary" gutterBottom>
          코로나19 발생 현황
        </Typography>
        <div className={classes.stat}>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.patients}</Typography>
              <b>치료중<br/>(격리중)</b>
          </div>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.cured}</Typography>
              <b>완치<br/>(격리해제)</b>
          </div>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.death}</Typography>
              <b>사망</b>
          </div>
          <div className={classes.statitem}>
              <Typography variant="h5" component="h2">{statData.confirmed}</Typography>
              <b>합계<br/>(확진)</b>
          </div>
      </div>
      <div className={classes.stat}>
          <div className={classes.statitem}>
            <span>검사중(검사진행) <b>{statData.checking}</b></span>
          </div>
          <div className={classes.statitem}>
            <span>결과음성 <b>{statData.resultNeg}</b></span>
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
          질병관리본부 최근 공지사항
          </Typography>
    <List>
      {newsData.map((item, i)=>(
        <div>
          <ListItem alignItems="flex-start" button
            onClick={()=>{
              window.open(item.link, "_blank")
              analytics.logEvent('link',{link: item.link});
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
        <Divider component="li" />
        </div>
      ))}
      </List>
      </CardContent>
      <CardActions>
        <Button size="small" href="http://ncov.mohw.go.kr/tcmBoardList.do?brdId=3">더보기</Button>
      </CardActions>
      </Card>
      <div className={classes.iconBtns}>
      <IconButton className={classes.iconBtnsItem}
        size="small" aria-label="close" color="inherit"
        href="https://youngbin.xyz">
        <LanguageIcon fontSize="large" />
      </IconButton>
      <IconButton className={classes.iconBtnsItem}
        size="small" aria-label="close" color="inherit"
        href="mailto:sukso96100@gmail.com">
        <AlternateEmailIcon fontSize="large" />
      </IconButton>
      <iframe className={classes.iconBtnsItem}
      src="https://ghbtns.com/github-btn.html?user=sukso96100&repo=covid19-push&type=star&count=true&size=large"
       frameborder="0" scrolling="0" width="160px" height="30px"></iframe>
      </div>
      <Snackbar
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        open={snackbar}
        autoHideDuration={6000}
        onClose={handleClose}
        message={snackMsg}
        action={
          <React.Fragment>
            <IconButton size="small" aria-label="close" color="inherit" onClick={handleClose}>
              <CloseIcon fontSize="small" />
            </IconButton>
          </React.Fragment>
        }
      />
  </div>
  );
}


async function tokenSaved(){
  let token = await localForage.getItem("token");
  return token != undefined && token === "";
}

