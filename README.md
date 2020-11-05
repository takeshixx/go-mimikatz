# go-mimikatz
A Go wrapper around a pre-compiled version of the Mimikatz executable for the purpose of anti-virus evasion.

### Requirements:
	MemoryModule => https://github.com/fancycode/MemoryModule

This application utilizes encryption to encrypt the main mimikatz binary. 

### Build Process:

The build process is pretty much completely automated in the Makefile. If you want to know more about how the build, 
take a look at the Makefile for more details:

```bash
make all
```

This command will create the final `bin/gogokatz.exe` file with the encrypted `mimikatz.exe` as an embedded resource.

*Note*: Everytime you run this command, the embedded `mimikatz.exe` will be encrypted with a different key.
