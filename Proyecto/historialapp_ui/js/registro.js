let apiRest = 'http://localhost:5000/historiasclinicas';

window.onload = () => {
    console.log("El app de las historias clínicas cargó");

    let nombre = document.getElementById('nombre');
    let edad = document.getElementById('edad');
    let peso = document.getElementById('peso');
    let registroFormulario = document.getElementById('registro');
    let estadoCivil = document.getElementById('estado-civil');
    let ocupacion = document.getElementById('ocupacion');
    let fumador = document.getElementById('fumador');
    let bebidasAlcoholicas = document.getElementById('bebidas-acoholicas');
    let medicado = document.getElementById('medicado');
    let hostMenuItems = document.querySelectorAll('.dropdown-item');
    let hostTitleView = document.getElementById('host-versions');

    /*JSON.stringify({
    			nombre: nombre.value,
    			edad: edad.value
    		})*/
    /* body: `nombre=${nombre.value}&edad=${edad.value}`*/
    registroFormulario.onsubmit = (event) => {
        event.preventDefault();
        console.log("Click en registrar");

        let historiaClinica = {
            nombre: nombre.value,
            edad: edad.value,
            weight: peso.value,
            is_smoker: fumador.checked,
            is_a_drug_user: medicado.checked,
            drinks_alcohol: bebidasAlcoholicas.checked,
            family_status: estadoCivil.value,
            ocupation: ocupacion.value
        };

        console.log(JSON.stringify(historiaClinica));

        fetch(apiRest, {
            method: 'POST',
            /*headers: {
      			'Content-Type': 'application/x-www-form-urlencoded'
    		},*/
            body: JSON.stringify(historiaClinica),
            mode: 'cors'
        }).then((response) => {
            if (response.status !== 200) {
                console.log("Parece que hay un problema. " + response.status);
                return;
            } else {
                window.location.href = "index.html";
            }

            /*response.json().then((data) => {
            	console.log(data);
            });*/
        }).catch((error) => {
            console.log("La petición falló", error);
        });
    };

    for (var i = 0; i < hostMenuItems.length; i++) {
        hostMenuItems[i].onclick = (event) => {
            event.preventDefault();
            let item = event.target;
            apiRest = item.getAttribute('value');

            for (var i = 0; i < hostMenuItems.length; i++) {
                hostMenuItems[i].classList.remove('active');
            }
            item.classList.add('active');
            hostTitleView.textContent = `Host: ${item.textContent} `;
        }
    }
};