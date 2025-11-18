# Feedback - flag matching problem using lvlath

This markdown describes every steps I had to do for integrating your script into OTO.


## Development process

1- Parameter struct

In the origin script, the FlagID is a simple string. But in OTO, the flag representation is my struct Parameter.
The FlagID will then be like a foreign key of my field Flag in the Parameter struct.

```go
type FlagID string
type Parameter struct {
	gorm.Model
	Flag          string         `gorm:"unique;not null"`
	Description   string         `gorm:"type:text"`
	BinID        int         	 `gorm:"not null"`
	Bin		  Binary
	RequiresRoot  bool           `gorm:"not null"`
	RequiresValue bool           `gorm:"not null"`
	ValueType     ValueType      `gorm:"not null"`
    ConflictsWith []*Parameter 	 `gorm:"many2many:flag_conflicts;joinForeignKey:flag_id;joinReferences:conflict_id"`
    DependsOn     []*Parameter 	 `gorm:"many2many:flag_dependencies;joinForeignKey:flag_id;joinReferences:depends_on_id"`
}
```

2- Schema struct integration

Now, I need to store the Schema struct somewhere. As the Schema is a set of every possible flags (in my case parameters), a schema is like a `Binary` which own the whole set of parameters.
In my case, the struct are representing my database schema. Which means my Parameter struct has a 

**First issue** : consistency

We need consistency for my requirement graph and conflict map. The issue is that core.Graph isn't serializable and the only solution for me is to duplicate the struct with a serializable version.

What I suggest is to add to your library serializable struct of your main data like Graph, Vertex, Edge and so on. It will create high level function that will allow instant storage of the struct.




## Conclusion

Overall, lvlath integrates well with OTO but there is still a bit to think of.
Here is a list of every points I came with :

- data serialization / deserialization : high level function to ease the persistency of the data in my database.
- semantic : a lot of features name are too math related and should be more friendly.
