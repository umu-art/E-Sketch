import logo from './logo.svg';
import './App.css';
import {UserApi} from 'nb_proxy_api';

let apiInstance = new UserApi();
apiInstance.apiClient.basePath = window.location.origin;
apiInstance.apiClient.defaultHeaders = {
    ...apiInstance.apiClient.defaultHeaders,
    'Authorization': 'Bearer token_example'
}
let opts = {
    'authDto': {
        'email': 'email_example',
        'password': 'password_example'
    }
};
apiInstance.login(opts).then(() => {
    console.log('API called successfully.');
}, (error) => {
    console.error(error);
});


function App() {
    return (
        <div className="App">
            <header className="App-header">
                <img src={logo} className="App-logo" alt="logo"/>
                <p>
                    Edit <code>src/App.js</code> and save to reload.
                </p>
                <a
                    className="App-link"
                    href="https://reactjs.org"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Learn React
                </a>
            </header>
        </div>
    );
}

export default App;
