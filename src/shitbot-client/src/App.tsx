import { useState } from 'react'
import './App.css'
import { TonConnectButton, TonConnectUIProvider } from '@tonconnect/ui-react'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <TonConnectUIProvider manifestUrl="https://ton-connect.github.io/demo-dapp-with-wallet/tonconnect-manifest.json">
        <div className="container" onClick={() => setCount((count) => count + 1)}>
          <span>My App with React UI</span>
          <TonConnectButton />
            count is {count}
        </div>
      </TonConnectUIProvider>
    </>
  )
}

export default App
