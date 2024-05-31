var client = null;

function startChat(channel) {
    client = new tmi.Client({
        channels: [channel]
    });
    client.connect();
    client.on('message', (channel, tags, message, self) => {
        let username = tags['display-name'];
        chatMessage(username, message);
    });
    client.on('join', (channel, username, self) => {
        if (self) {
            onSelfJoinMessage();
        }
    });
}

function stopChat() {
    client.disconnect();
    client = null;
}
