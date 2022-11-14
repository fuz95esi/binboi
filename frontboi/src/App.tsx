import React from 'react';
import './App.css';
import { CodeTxt } from './components/CodeTxt'
import { FormHero } from './components/FormHero'

function App() {
  return (
    <div className="App">
      <p style={{width: "40em", margin: "0.2em"}}>
        👋 I'm <CodeTxt>binboi</CodeTxt> and I can generate a full year's set of bin collection reminders for your calendar. 
      </p>
      <p style={{width: "40em", margin: "0.2em", fontWeight: "600"}}>
        Enter your postcode to get started.
      </p>
      <p style={{fontStyle: "italic", fontSize: "calc(4px + 2vmin)", margin: "1em"}}>
        (I only working for Reading council at the moment.)
      </p>
      <FormHero/>
    </div>
  );
}

export default App;
