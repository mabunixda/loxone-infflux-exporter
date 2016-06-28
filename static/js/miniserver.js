
function loadLoxConfig() {
    var tmp = localStorage.getItem("LoxoneConfig");
    var loxConfig = null;
    if(tmp) {
        loxConfig = JSON.parse(tmp);
    }
    var remote ="/loxapp3.json" 
    $.getJSON(remote, function(loxapp) {
        if ( loxConfig == null || loxConfig.lastModified < loxapp.lastModified ) {
            localStorage.setItem("LoxoneConfig", JSON.stringify(loxapp));
        }   
    });
}