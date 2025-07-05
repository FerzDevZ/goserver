// reload.js
(function() {
  var ws = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/ws');
  ws.onmessage = function(msg) {
    if (msg.data === 'reload') {
      window.location.reload();
    }
  };
  ws.onclose = function() {
    console.warn('Live reload connection closed.');
  };
})();
