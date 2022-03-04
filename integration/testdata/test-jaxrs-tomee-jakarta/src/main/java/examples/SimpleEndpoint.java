package examples;

import jakarta.enterprise.context.RequestScoped;
import jakarta.json.Json;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/")
@RequestScoped
public class SimpleEndpoint {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public String health() {
        return Json.createObjectBuilder().add("application_status", "UP").build().toString();
    }

}
