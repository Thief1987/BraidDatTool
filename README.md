Usage : 
  
    unpack: -u archive_name 
    repack: -r archive name compression_level

When repacking, optionally you can specify compression level, legit values are from -4(fastest) to 9(slowest).
Default value is 6(devs used it), but it's pretty slow, very slow I would say, so I decided to add this option at least for the testing purposes.
If not specified or value is invalid< it will be set to 6.
