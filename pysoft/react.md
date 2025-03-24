# Detección de cambios

React trabaja en un SPA (Single Page Aplication), hay un index.html y un script que carga todo lo necesario.

## Trigger 

Sirven para detectar cambios en la página. Es un evento que va a iniciar un evento de **render**.

Por ejemplo, un botón, un estado, una llamada asíncrona, una conexión a una API.

### Trigger inicial

Monta el componente inicial.

### Re-Render

React trabaja con DOM y DOM-virtual. En el DOM está el contenido ya montado de la página (el render html) y el DOM-virtual se actualiza al haber un trigger, si el render html cambió se recarga la página.

Las cosas que hacen un Render pueden ser:

1. Mount -> montado de un nuevo componente
2. Cambio de estado
3. Async (aunque también da un cambio de estado).


## Hooks

Son componentes de React para producir cambios de estado.

### `useState()`

Para hacer una conexión entre el render con lo que es el contenido de la función estamos hablando de un estado. Una variable local no dispara ningún trigger al modificarla.

```js

// Si pulso el botón no funcionaría

let localCount = 0;

const updateLocalCount = () => {
    setCount((localCount) => localCount + 1)
}


<Button onClick={updateLocalCount} />

<label>{localCount}</label>

```

```js

// Esto si que funciona al pulsar el botón

const [count, setCount] = useState(0)

const updateCount = () => {
    setCount((count) => count + 1)
}

<Button onClick={updateCount} />

<label>{count}</label>

``` 


### `useEffect()`

**Permite sincronizar con entidades externas, por ejemplo, para conectarse a un endpoint**.

El contenido de un useEffect() se ejecuta al producirse un render o cuando hay una modificación sobre un useState incluido en la lista de dependencias.

Un `return` dentro de useEffect(), permite liberar memoria manualmente, puede ser útil para finalizar funciones asíncronas.

Se puede usar más de un useEffect() a la vez.

> **USAR PARA SINCRONIZAR CON COSAS EXTERNAS**

## Llamadas a endpoits

Deben ser funciones asíncronas ya que pueden llevar mucho tiempo. Por ejemplo:

```js
const [data, setData] = useState([])
const [loading, setLoading] = useState(false)
const [error, setError] = useState("")

const fetchData = async () => {
    setLoading(true)
    try {

        const response = await fetch("https://api.example.com/data")

        if(!response.ok) {
            throw new Error("Error al obtener datos")
        }

        const jsonData = await response.json()

        setData(jsonData)

    } catch (err) {
        setError(err as string)
    } finally {
        setLoading(false)
    }
}


if(loading){
    return <div>Cargando...</div>
}

if(error) {
    return <div>UPS!!!!!</div>
}
```

En el código anterior habría que tener cuidado con bucles infinitos. Si añadiesemos el código de aquí abajo, fetchData() estaría ejecutándose todo el rato.

```js

useEffect(() => {
    fetchData()
}, [data]) 

```

Para solucionarlo se puede usar `useCallback()`.

### Custom hooks

Se pueden crear custom hooks: por ejemplo, en un fichero `/src/hooks/useFetch.js`.

```js

type Data<T> = T | null;
type ErrorType = Error | null;

interface Params<T> {
    data: Data<T>;
    loading: boolean;
    error: ErrorType;
}

export const useFetch = <T>(url: string): Params<T> => {
    const [data, setData] = useState<Data<T>>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<ErrorType>(null)

    useEffect(() => {
        const controller = new AbortController();


        const fetchData = async () => {
            setLoading(true)
            try{
                const response = await fetch(url, controller);
                if (!response.ok) {
                    throw new Error("Error en la petición")
                }

                const jsonData: T = await response.json();
                setData(jsonData)
                setError(null)
            } catch (err) {
                setError(err as Error)
            } finally {
                setLoading(false)
            }
        }

        fetchData();

        return () => {
            controller.abort();
        }

        
    }, [url])
 
    return {data, loading, error}
}


```

