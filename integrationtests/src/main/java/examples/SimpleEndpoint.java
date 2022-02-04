package examples;

import javax.enterprise.context.RequestScoped;
import javax.json.Json;
import javax.json.JsonArrayBuilder;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

@Path("/")
@RequestScoped
public class SimpleEndpoint {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public String getTopCDs() {

        final JsonArrayBuilder array = Json.createArrayBuilder();
        final List<Integer> randomCDs = getRandomNumbers();
        for (final Integer randomCD : randomCDs) {
            array.add(Json.createObjectBuilder().add("id", randomCD));
        }
        return array.build().toString();
    }

    private List<Integer> getRandomNumbers() {
        final List<Integer> randomCDs = new ArrayList<>();
        final Random r = new Random();
        randomCDs.add(r.nextInt(100) + 1101);
        randomCDs.add(r.nextInt(100) + 1101);
        randomCDs.add(r.nextInt(100) + 1101);
        randomCDs.add(r.nextInt(100) + 1101);
        randomCDs.add(r.nextInt(100) + 1101);

        return randomCDs;
    }
}
