window.addEventListener("load", function(){
    var form2 = document.getElementById("form2")
    connectBut = this.document.getElementById("connect")
    loadh()
    connectBut.addEventListener("mousemove", function() {
        connectBut.style.display = "none"
        var socket = new WebSocket("ws://127.0.0.1:6060/ws");
        socket.onopen = function() {

            socket.send(document.getElementById("name").value)
            socket.send(document.getElementById("pass").value)
            setTimeout(function(){
              document.getElementById("former").style.display = "block"

            }, 500) 
          };
          
          socket.onclose = function(event) {
            if (event.wasClean) {
                this.console.log('Closed clean');
            } else {
                this.console.log('Connection tear');
            }
            this.console.log('Error: ' + event.code + '. Reason: ' + event.reason);
            connectBut.style.display = "inline"
            document.getElementById("name").style.display = "inline"
            document.getElementById("pass").style.display = "inline"
            document.getElementById("former").style.display = "none"
            document.getElementById("answer").innerHTML = ""
          };
          
          socket.onmessage = function(event) {
            if (event.data == "wrong_reg") {
              document.getElementById("former").style.display = "none"
              alert("Wrong password")
              socket.close()
            }
            let answer = document.getElementById("answer")
            answer.innerHTML = event.data
          };
          
          socket.onerror = function(error) {
            alert("Ошибка " + error.message);
          };
          let command = document.getElementById("command")
          command.value = ""
          form2.addEventListener("submit", function(event){
            event.preventDefault();
            if (command.value == "exit") {
              socket.close()
            }
            socket.send(command.value)
          })
    })
})

