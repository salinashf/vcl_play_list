
# Play List VLC

Esta aplicación te permite generar el archivo de **PlayList** de cursos o canciones, este  genera un archivo **XML**  del contenido del curso o canciones.

Este archivo hay que abrirlo con VLC y  carga automáticamente todo el curso ... para evitarte abrir de uno en uno los videos


Estructura del XML

```xml
<playlist xmlns="http://xspf.org/ns/0/" xmlns:vlc="http://www.videolan.org/vlc/playlist/ns/0/" version="1">
	<title>PHP y Composer</title>
	<trackList>
		<track>
			<location>file:///A:/Cursos/Curso PHP/Introduccion/Qué aprende sobre PHP </location>
			<title>Qué aprenderás sobre PHP </title>
			<album>Introduccion</album>
			<trackNum>1</trackNum>
			<extension application="http://www.videolan.org/vlc/playlist/0">
				<vlc:id>1</vlc:id>
			</extension>
		</track>
		<track>
			<location>file:///A:/Cursos/Curso PHP/Introduccion/Qué aprendera sobre Composer</location>
			<title>Qué aprenderás sobre PHP </title>
			<album>Introduccion</album>
			<trackNum>2</trackNum>
			<extension application="http://www.videolan.org/vlc/playlist/0">
				<vlc:id>2</vlc:id>
			</extension>
		</track>
	</trackList>
</playlist>

```
# Como Generar

Para generar debe de pasar la ruta con el prefijo de argumento **-path**

> **play_list.exe -path** *"A:\\Cursos\\Curso PHP\\"* 

## Documentacion del formato 
Para mayor información del formato **XSPF** , link de Wikipedia 
https://en.wikipedia.org/wiki/XML_Shareable_Playlist_Format
