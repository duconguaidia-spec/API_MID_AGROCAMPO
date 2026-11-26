param(
    [string]$BaseUrl = 'http://localhost:4200',
    [string]$CrudUrl = 'http://localhost:8087',
    [switch]$Escritura,
    [pscredential]$Credencial
)

$ErrorActionPreference = 'Stop'
$api = $BaseUrl.TrimEnd('/') + '/api/v1/veterinarias'
$crud = $CrudUrl.TrimEnd('/') + '/v1'
$sesion = [Microsoft.PowerShell.Commands.WebRequestSession]::new()

function Invoke-Api([string]$Method, [string]$Uri, $Data = $null) {
    $parameters = @{ Method = $Method; Uri = $Uri; TimeoutSec = 30; WebSession = $sesion }
    if ($null -ne $Data) {
        $parameters.ContentType = 'application/json; charset=utf-8'
        $parameters.Body = [System.Text.Encoding]::UTF8.GetBytes(($Data | ConvertTo-Json -Depth 15))
    }
    $response = Invoke-RestMethod @parameters
    if ($response.success -ne $true) { throw "La API rechazó $Method $Uri" }
    return $response
}

$list = Invoke-Api GET $api
Write-Output "Listado correcto: $(@($list.data).Count) veterinarias."
if (-not $Escritura) { return }

if (-not $Credencial) { $Credencial = Get-Credential -Message 'Cuenta administradora o veterinaria de AgroCampo' }
if (-not $Credencial) { throw 'La prueba de edición requiere una cuenta de AgroCampo.' }
$cuenta = Invoke-Api POST ($BaseUrl.TrimEnd('/') + '/api/v1/cuenta/login') @{
    correo = $Credencial.UserName
    contrasena = $Credencial.GetNetworkCredential().Password
}
if ($cuenta.data.rol -notin @('administrador', 'veterinario')) { throw 'La cuenta no tiene permiso de edición.' }

# Solo se eliminan los registros creados por esta ejecución.
$name = 'Prueba conexión ' + [guid]::NewGuid().ToString('N')
$vetId = 0
$productIds = [System.Collections.Generic.List[int]]::new()
$serviceIds = [System.Collections.Generic.List[int]]::new()
try {
    $created = Invoke-Api POST $api @{
        nombre = $name; descripcion = 'Registro temporal de integración'; direccion = 'Calle de prueba 1'
        ciudad = 'Yopal'; telefono = '3000000000'; horario_atencion = '8:00 a 17:00'
        productos = @(@{ nombre = 'Producto de prueba'; descripcion = 'Temporal'; precio = 1000 })
        servicios = @(@{ nombre = 'Consulta de prueba'; descripcion = 'Temporal'; precio = 2000 })
    }
    $vetId = [int]$created.data.veterinaria.id
    if ($vetId -le 0) { throw 'Falta el ID creado por PostgreSQL.' }
    foreach ($product in $created.data.productos) { $productIds.Add([int]$product.id) }
    foreach ($service in $created.data.servicios) { $serviceIds.Add([int]$service.id) }

    $detail = Invoke-Api GET "$api/$vetId"
    if (@($detail.data.productos).Count -ne 1 -or @($detail.data.servicios).Count -ne 1) { throw 'El detalle no contiene el producto y servicio creados.' }

    Invoke-Api PUT "$api/$vetId" @{ descripcion = 'Descripción actualizada'; horario_atencion = '9:00 a 18:00' } | Out-Null
    $detail = Invoke-Api GET "$api/$vetId"
    if ($detail.data.veterinaria.descripcion -ne 'Descripción actualizada' -or $detail.data.veterinaria.ciudad -ne 'Yopal') { throw 'La edición no se conservó o perdió campos.' }

    $productId = $productIds[0]
    Invoke-Api PUT "$api/$vetId/productos/$productId" @{ nombre = 'Producto actualizado'; descripcion = 'Temporal'; precio = 1500 } | Out-Null
    $newProduct = Invoke-Api POST "$api/$vetId/productos" @{ nombre = 'Otro producto'; precio = 500 }
    $newProductId = [int]$newProduct.data.id
    if ($newProductId -le 0) { throw 'Falta el ID del nuevo producto.' }
    $productIds.Add($newProductId)
    $detail = Invoke-Api GET "$api/$vetId"
    $updatedProduct = @($detail.data.productos | Where-Object { $_.id -eq $productId })
    if ($updatedProduct.Count -ne 1 -or $updatedProduct[0].nombre -ne 'Producto actualizado' -or $updatedProduct[0].precio -ne 1500) { throw 'La edición del producto no persistió.' }
    if (@($detail.data.productos).Count -ne 2) { throw 'No se guardó el segundo producto.' }

    Invoke-Api DELETE "$api/$vetId/productos/$newProductId" | Out-Null
    $productIds.Remove($newProductId) | Out-Null
    $detail = Invoke-Api GET "$api/$vetId"
    if (@($detail.data.productos).Count -ne 1) { throw 'El producto eliminado sigue apareciendo.' }
    Write-Output 'Lectura, registro, edición y eliminación de productos: correctos y persistidos en PostgreSQL.'
}
finally {
    if ($vetId -gt 0) {
        $original = Invoke-Api GET "$crud/veterinaria/$vetId"
        if ($original.data.nombre -ne $name) { throw 'No se limpia: el registro no coincide con esta ejecución.' }
        foreach ($productId in $productIds) { Invoke-Api DELETE "$crud/producto_veterinaria/$productId" | Out-Null }
        foreach ($serviceId in $serviceIds) { Invoke-Api DELETE "$crud/servicioveterinaria/$serviceId" | Out-Null }
        Invoke-Api DELETE "$crud/veterinaria/$vetId" | Out-Null
        Write-Output 'Registros temporales eliminados.'
    }
}
