import './App.css';
import React from 'react';
import { Chat as ApiContract, Message } from '../../client-ts/client.gen';
import { flushSync } from 'react-dom';

const customFetch = window.fetch.bind(window);
const contract = new ApiContract('http://localhost:4242', customFetch);

function Chat() {
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [connected, setConnected] = React.useState(false);
  const [userMessage, setUserMessage] = React.useState('');
  const chatRef = React.useRef<HTMLOListElement>(null);

  React.useEffect(() => {
    const { subscribe } = contract.subscribeMessages();

    const unsubscribe = subscribe({
      onData: (data) => {
        flushSync(() => {
          setMessages((prevMessages) => [...prevMessages, data]);
          const lastMessage = chatRef.current?.lastElementChild;
          lastMessage?.scrollIntoView({
            block: 'end',
            behavior: 'smooth',
            inline: 'nearest',
          });
        });
      },
      onOpen: () => {
        setConnected(true);
      },
      onClose: () => {
        setConnected(false);
      },
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const onUserMessageChanged = (ev: React.ChangeEvent<HTMLInputElement>) => {
    setUserMessage(ev.currentTarget.value);
  };

  const sendMessage = (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    contract
      .sendMessage({ author: 'FE', msg: userMessage })
      .then(() => setUserMessage(''));
  };

  return (
    <div className="chat">
      <h1>Chat {connected ? '✅' : '❌'}</h1>
      <ol className="messages" ref={chatRef}>
        {messages.map((message) => (
          <li key={message.id}>
            {message.author} - {message.msg}
          </li>
        ))}
      </ol>
      <form onSubmit={sendMessage}>
        <input value={userMessage} onChange={onUserMessageChanged} />
        <button type="submit">Send</button>
      </form>
    </div>
  );
}

function App() {
  const [show, setShow] = React.useState(false);
  const [username, setUsername] = React.useState('');

  return (
    <div className="App">
      <button onClick={() => setShow((prev) => !prev)}>
        {show ? 'Disconnect' : 'Connect'}
      </button>
      <input value={username} />
      {show ? <Chat /> : null}
    </div>
  );
}

export default App;
