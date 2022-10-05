import React from 'react';
import { flushSync } from 'react-dom';
import { Chat as ApiContract, Message } from '../../client-ts/client.gen';
import './App.css';

const customFetch = window.fetch.bind(window);
const contract = new ApiContract('http://localhost:4242', customFetch);

type Props = {
  username: string;
};

function Chat(props: Props) {
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [connected, setConnected] = React.useState(false);
  const [userMessage, setUserMessage] = React.useState('');
  const chatRef = React.useRef<HTMLOListElement>(null);
  const { username } = props;

  React.useEffect(() => {
    const { subscribe } = contract.subscribeMessages();

    const unsubscribe = subscribe({
      onData: (data) => {
        flushSync(() => {
          setMessages((prevMessages) => [...prevMessages, data]);
        });
        const lastMessage = chatRef.current?.lastElementChild;
        lastMessage?.scrollIntoView({
          block: 'end',
          behavior: 'smooth',
          inline: 'nearest',
        });
      },
      onError: (e) => {
        console.error(e);
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
      .sendMessage({ author: username, msg: userMessage })
      .then(() => setUserMessage(''));
  };

  return (
    <div className="chat">
      <h2>Connected {connected ? '✅' : '❌'}</h2>
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
  const [username, setUsername] = React.useState('Your name');

  return (
    <div className="App">
      <input value={username} onChange={(e) => setUsername(e.target.value)} />
      <button onClick={() => setShow((prev) => !prev)}>
        {show ? 'Disconnect' : 'Connect'}
      </button>
      {show ? <Chat username={username} /> : null}
    </div>
  );
}

export default App;
