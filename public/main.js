$(function() {
    if (!window["WebSocket"]) {
        return;
    }

    var content = $("#content");
    var conn = new WebSocket('ws://localhost:8080/ws');

    // Textarea is editable only when socket is opened.
    conn.onopen = function(e) {
        content.attr("disabled", false);

        // Подписываемся на канал user_1
        conn.send('{"method": "subscribe", "params": {"channel": "user_1"}}');
    };

    conn.onclose = function(e) {
        content.attr("disabled", true);
    };

    conn.onerror = function(error) {
        console.log(error)
    };

    // Whenever we receive a message, update textarea
    conn.onmessage = function(e) {
        if (e.data != content.val()) {
            content.val(e.data);
        }
    };
});