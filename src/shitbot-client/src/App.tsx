import { useEffect, useState } from 'react'
import './App.css'
import { TonConnectButton, TonConnectUIProvider } from '@tonconnect/ui-react'
import { Button } from '@mui/material';
import WebApp from '@twa-dev/sdk';

function App() {
  const [count, setCount] = useState(0);
  const [user, setUser] = useState<any>({});
  const [stat, setStat] = useState<any>({});

  useEffect(() => {
    const fetchUser = async() => {
      let response = await fetch(`http:127.0.0.1/users?${WebApp.initData}`) ;
      let data: any = await response.json();
      if (data.success) {
        setUser(data.data);
        setCount(data.data.coins);
        setStat(data.data.stat);
      } else if (response.status === 404) {
        response = await fetch(`http:127.0.0.1/users?${WebApp.initData}`, {
          method: "POST"
        });
        data = await response.json();
        if (data.success) {
          setUser(data.data);
          setCount(data.data.coins);
          setStat(data.data.stat);
        } else {
          console.log(data.message)
        }
      }
    }
    fetchUser()
    .catch((err) => {
      console.log(err);
    });
 }, []);
 
  return (
    <>
      <TonConnectUIProvider manifestUrl="https://ton-connect.github.io/demo-dapp-with-wallet/tonconnect-manifest.json">
        <div className="flex flex-col items-center mt-12">
          <TonConnectButton/>
          <div className="mt-8"></div>
          <Button variant="text" onClick={() => setCount(count => count + 1)}>Click</Button>
            {count}
          <div className='flex flex-col items-center mt-8'>
            <span  className='mt-4'>{stat.earned}</span>
            <span  className='mt-4'>{`https://t.me/dolbikshit_bot?start=${user.referralCode}`}</span>

          </div>
        </div>
      </TonConnectUIProvider>
    </>
  )
}

export default App
