# GOPOLY

-g
'AdvertBase
subtypes=RentAdvert,SellAdvert
marker_method=i{{.Type}}
discriminator.field=type
discriminator.mapping=rent:RentAdvert,sell:SellAdvert
decoding_strategy=discriminator
template_file="template.txt"
filename=out.gen.go';