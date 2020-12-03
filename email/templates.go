package email

const addressConfirmationTemplate = `
<style>
    body {
        text-align: center;
        width: 40%;
        margin: auto;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    }

    h1 {
        color: red;
    }

    button {
        background-color: #FD5D5D;
        border: none;
        padding: 10px 20px;
        color: white;
        font-weight: 600;
        border-radius: 10px;
        margin: auto;
    }

    .button-div {
        width: 100%;
        text-align: center;
    }

    .logo {
        width: 80%;
    }
</style>

<body>
    <h1>Hej {{.Name}}</h1>
    <p>
        Mange tak for at du har oprettet en bruger hos Landlyst kro og hotel!
    </p>
    <p>
        Vi ser frem til mange besøg, og vi glæder os til at se dig.
    </p>
    <p>For at bekræfte din profil på landlysthotel.dk, skal du blot <strong>klikke på linket nedenunder</strong>.</p>
    <br>
    <div class="button-div">
        <button>Bekræft email</button>
        
    </div>
    <br>
    Med venlig hilsen<br />
    <br>
    <img class="logo" src="../public/images/Logo@2x.png"></img>
</body>
`
