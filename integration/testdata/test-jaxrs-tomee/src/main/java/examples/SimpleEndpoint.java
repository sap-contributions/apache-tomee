package examples;

import javax.enterprise.context.RequestScoped;
import javax.json.Json;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("/")
@RequestScoped
public class SimpleEndpoint {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public String health() {
        return Json.createObjectBuilder().add("application_status", "UP").build().toString();
    }

}
