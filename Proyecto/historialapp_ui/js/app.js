let apiRest = 'http://localhost:5000/historiasclinicas';

window.onload = () => {
	console.log("El app de las historias clínicas cargó");

	let historiasView = document.getElementById('historias');
	let obtenerBoton = document.getElementById('obtener');
	let hostMenuItems = document.querySelectorAll('.dropdown-item');
	let hostTitleView = document.getElementById('host-versions');

	function obtenerHistoriasClinicas() {
		fetch(apiRest, {mode: 'cors'}).then((response) => {
			return response.json();
		}).then((data) => {
			console.log(data);
			historiasView.innerHTML = '';
			let contador = 0;

			for (var i = 0; i < data.length; i++) {
				contador++;
				let h = data[i];

				let fila = document.createElement("TR");
				let columnaIndex = document.createElement("TH");
				columnaIndex.setAttribute("scope", "row");
				columnaIndex.textContent = contador;

				let columnaNombre = document.createElement("TD");
				columnaNombre.textContent = h.nombre;

				let columnaEdad = document.createElement("TD");
				columnaEdad.textContent = h.edad;

				let columnaPeso = document.createElement("TD");
				columnaPeso.textContent = `${h.weight} Kg.`;

				let columnaFumador = document.createElement("TD");
				columnaFumador.textContent = h.is_smoker ? 'Sí' : 'No';

				let columnaMedicado = document.createElement("TD");
				columnaMedicado.textContent = h.is_a_drug_user ? 'Sí' : 'No';

				let columnaBebidasAlcoholicas = document.createElement("TD");
				columnaBebidasAlcoholicas.textContent = h.drinks_alcohol ? 'Sí' : 'No';

				let columnaEstadoCivil = document.createElement("TD");
				columnaEstadoCivil.textContent = h.family_status;

				let columnaOcupacion = document.createElement("TD");
				columnaOcupacion.textContent = h.ocupation;							

				fila.appendChild(columnaIndex);
				fila.appendChild(columnaNombre);
				fila.appendChild(columnaEdad);
				fila.appendChild(columnaPeso);
				fila.appendChild(columnaFumador);
				fila.appendChild(columnaMedicado);
				fila.appendChild(columnaBebidasAlcoholicas);
				fila.appendChild(columnaEstadoCivil);
				fila.appendChild(columnaOcupacion);

				historiasView.appendChild(fila);
			}
			//historiasView.textContent = text;
		}).catch((error) => {
			console.log("La petición falló", error);
		});		
	}

	obtenerBoton.onclick = (event) => {
		event.preventDefault();
		console.log("Click en obtener");
		obtenerHistoriasClinicas();
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

	//obtenerHistoriasClinicas();
};